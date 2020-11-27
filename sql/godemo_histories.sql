CREATE DATABASE  IF NOT EXISTS `godemo` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `godemo`;
-- MySQL dump 10.13  Distrib 8.0.22, for macos10.15 (x86_64)
--
-- Host: localhost    Database: godemo
-- ------------------------------------------------------
-- Server version	8.0.22

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `histories`
--

DROP TABLE IF EXISTS `histories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `histories` (
  `history_id` int unsigned NOT NULL AUTO_INCREMENT,
  `gi_id` int unsigned DEFAULT NULL,
  `course_time` datetime DEFAULT NULL,
  `course_name` varchar(255) DEFAULT NULL,
  `course_score` double DEFAULT NULL,
  PRIMARY KEY (`history_id`)
) ENGINE=InnoDB AUTO_INCREMENT=103 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `histories`
--

LOCK TABLES `histories` WRITE;
/*!40000 ALTER TABLE `histories` DISABLE KEYS */;
INSERT INTO `histories` VALUES (1,1,'2020-11-21 01:00:00','something new1',3.2),(2,1,'2020-11-21 01:00:00','Super Course',4.2),(3,1,'2020-11-21 01:00:00','new insert name',1.2),(4,1,'2020-11-21 01:00:00','Super . asdasd',2.3),(5,1,'2020-11-21 01:00:00','some course',1.1),(6,1,'2020-11-21 01:00:00','some course',3.11),(7,1,'2020-11-21 01:00:00','some course',1.5),(8,1,'2020-11-21 01:00:00','some course',6.3),(9,1,'2020-11-21 01:00:00','some course',4),(10,2,'2020-11-21 01:00:00','some course',6.1),(11,2,'2020-11-21 01:00:00','some course',0.2),(12,2,'2020-11-21 01:00:00','some course',1.33),(13,2,'2020-11-21 01:00:00','some course',4.5),(14,2,'2020-11-21 01:00:00','some course',3.2),(15,2,'2020-11-21 01:00:00','some course',3.11),(16,3,'2020-11-21 01:00:00','123 things 123',5.66),(17,3,'2020-11-21 01:00:00','some course',1.54),(18,3,'2020-11-21 01:00:00','Something',99.3),(19,3,'2020-11-21 01:00:00','Something',929.3),(20,3,'2020-11-21 01:00:00','Something',929.3),(21,3,'2020-11-21 01:00:00','Something',929.3);
/*!40000 ALTER TABLE `histories` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-11-27 10:18:40
