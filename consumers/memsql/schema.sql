DROP DATABASE IF EXISTS vflow;
CREATE DATABASE vflow;
USE vflow;

CREATE TABLE `samples` (
  `device` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `src` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `dst` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `nextHop` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `srcASN` int(11) UNSIGNED DEFAULT NULL,
  `dstASN` int(11) UNSIGNED DEFAULT NULL,
  `proto` int(11) DEFAULT NULL,
  `srcPort` int(11) DEFAULT NULL,
  `dstPort` int(11) DEFAULT NULL,
  `tcpFlags` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `ingressIf` int(11) DEFAULT NULL,
  `egressIf` int(11) DEFAULT NULL,
  `bytes` int(11) DEFAULT NULL,
  `datetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
  /*!90618 , SHARD KEY () */ 
)
