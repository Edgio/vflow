Name:       vflow
Version:    %VERSION%
Release:    0
Group:      Application
URL:        https://github.com/VerizonDigital/vflow
License:    Apache-2
Summary:    IPFIX/sFlow/Netflow collector
Source0:    vflow
Source1:    vflow.conf
Source2:    vflow.service

%description
High-performance, scalable and reliable IPFIX, sFlow and Netflow collector

%prep

%install
rm -rf %{buildroot}

mkdir -p %{buildroot}/usr/bin
mkdir -p %{buildroot}/usr/local/vflow
mkdir -p %{buildroot}/etc/vflow
mkdir -p %{buildroot}/etc/init.d
cp -Rf %{SOURCE0} %{buildroot}/usr/bin/
cp -Rf %{SOURCE1} %{buildroot}/etc/vflow/
cp -Rf %{SOURCE2} %{buildroot}/etc/init.d/vflow

%files
/usr/bin/vflow
/etc/vflow/vflow.conf
/etc/init.d/vflow
