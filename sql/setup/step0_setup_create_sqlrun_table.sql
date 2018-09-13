CREATE TABLE IF NOT EXISTS `SqlRunner` (
  `SqlRunnerId` int(11) NOT NULL AUTO_INCREMENT,
  `SqlRunnerName` varchar(2048) DEFAULT NULL,
  `StatusId` int(11) NOT NULL,
  `ErrorCount` int(11) NOT NULL DEFAULT 0,
  `ErrorText` text DEFAULT NULL,
  `LastCreated` timestamp(3) NOT NULL DEFAULT current_timestamp(3),
  `LastModified` timestamp(3) NOT NULL DEFAULT current_timestamp(3) ON UPDATE current_timestamp(3),
  PRIMARY KEY (`SqlRunnerId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;