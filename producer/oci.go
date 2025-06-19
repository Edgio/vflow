//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    oci.go
//: details: vflow OCI Object Store producer
//: author:  Enhanced by Claude Code
//: date:    2024-12-18
//:
//: Licensed under the Apache License, Version 2.0 (the "License");
//: you may not use this file except in compliance with the License.
//: You may obtain a copy of the License at
//:
//:     http://www.apache.org/licenses/LICENSE-2.0
//:
//: Unless required by applicable law or agreed to in writing, software
//: distributed under the License is distributed on an "AS IS" BASIS,
//: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//: See the License for the specific language governing permissions and
//: limitations under the License.
//: ----------------------------------------------------------------------------

package producer

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v2"
)

// OCIConfig represents OCI Object Store configuration
type OCIConfig struct {
	Region         string `yaml:"region"`
	TenancyOCID    string `yaml:"tenancy_ocid"`
	UserOCID       string `yaml:"user_ocid"`
	Fingerprint    string `yaml:"fingerprint"`
	PrivateKeyPath string `yaml:"private_key_path"`
	BucketName     string `yaml:"bucket_name"`
	Namespace      string `yaml:"namespace"`
	ChunkSizeMB    int    `yaml:"chunk_size_mb"`
	FlushInterval  string `yaml:"flush_interval"`
}

// OCI represents OCI Object Store producer
type OCI struct {
	config        OCIConfig
	logger        *log.Logger
	buffer        *bytes.Buffer
	bufferSize    int64
	maxChunkSize  int64
	flushInterval time.Duration
	fileCounter   int64
}

// setup initializes the OCI producer
func (o *OCI) setup(configFile string, logger *log.Logger) error {
	o.logger = logger
	o.buffer = &bytes.Buffer{}
	o.bufferSize = 0
	o.fileCounter = 0

	// Load configuration
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read OCI config file %s: %v", configFile, err)
	}

	err = yaml.Unmarshal(configData, &o.config)
	if err != nil {
		return fmt.Errorf("failed to parse OCI config: %v", err)
	}

	// Set defaults
	if o.config.ChunkSizeMB == 0 {
		o.config.ChunkSizeMB = 10
	}
	if o.config.FlushInterval == "" {
		o.config.FlushInterval = "60s"
	}

	o.maxChunkSize = int64(o.config.ChunkSizeMB) * 1024 * 1024

	o.flushInterval, err = time.ParseDuration(o.config.FlushInterval)
	if err != nil {
		return fmt.Errorf("invalid flush interval: %v", err)
	}

	// Validate required configuration
	if o.config.BucketName == "" || o.config.Namespace == "" {
		return fmt.Errorf("bucket_name and namespace are required in OCI config")
	}

	o.logger.Printf("OCI producer initialized: bucket=%s, namespace=%s, chunk_size=%dMB, flush_interval=%s",
		o.config.BucketName, o.config.Namespace, o.config.ChunkSizeMB, o.config.FlushInterval)

	return nil
}

// inputMsg processes messages from the channel and uploads to OCI
func (o *OCI) inputMsg(topic string, ch chan []byte, errorCount *uint64) {
	ticker := time.NewTicker(o.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				// Channel closed, flush remaining data and exit
				if o.buffer.Len() > 0 {
					o.flushToOCI()
				}
				return
			}

			// Add message to buffer
			o.addToBuffer(msg)

			// Check if we need to flush based on size
			if o.bufferSize >= o.maxChunkSize {
				o.flushToOCI()
			}

		case <-ticker.C:
			// Flush on timer if buffer has data
			if o.buffer.Len() > 0 {
				o.flushToOCI()
			}
		}
	}
}

// addToBuffer adds a JSON message to the current buffer
func (o *OCI) addToBuffer(msg []byte) {
	// Add newline separator for JSON records
	o.buffer.Write(msg)
	o.buffer.WriteByte('\n')
	o.bufferSize += int64(len(msg) + 1)
}

// flushToOCI uploads the current buffer to OCI Object Store
func (o *OCI) flushToOCI() {
	if o.buffer.Len() == 0 {
		return
	}

	// Generate filename with timestamp and counter
	timestamp := time.Now().UTC().Format("20060102_150405")
	fileCounter := atomic.AddInt64(&o.fileCounter, 1)
	filename := fmt.Sprintf("sflow_%s_%03d.json", timestamp, fileCounter)

	// For now, write to local file as a placeholder
	// In production, this would upload to OCI Object Store using OCI SDK
	err := o.writeToLocalFile(filename, o.buffer.Bytes())
	if err != nil {
		o.logger.Printf("Error writing sFlow data to file %s: %v", filename, err)
		return
	}

	o.logger.Printf("sFlow data written to file: %s (size: %d bytes)", filename, o.buffer.Len())

	// Reset buffer
	o.buffer.Reset()
	o.bufferSize = 0
}

// writeToLocalFile writes data to a local file (placeholder for OCI upload)
func (o *OCI) writeToLocalFile(filename string, data []byte) error {
	// Create directory if it doesn't exist
	dir := "/tmp/vflow-oci"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write file
	filepath := fmt.Sprintf("%s/%s", dir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return writer.Flush()
}

// TODO: Implement actual OCI SDK upload
// This function would replace writeToLocalFile in production
func (o *OCI) uploadToOCI(filename string, data []byte) error {
	// Implementation would use OCI SDK:
	// 1. Create OCI client with credentials
	// 2. Upload data to bucket using PutObject
	// 3. Handle retries and error cases
	return fmt.Errorf("OCI SDK upload not implemented yet")
}