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
-- Table structure for table `permissions`
--

DROP TABLE IF EXISTS `permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `permissions` (
  `perm_id` int unsigned NOT NULL AUTO_INCREMENT,
  `column_id` int unsigned DEFAULT NULL,
  `group_id` int unsigned DEFAULT NULL,
  `perm_mask` int unsigned DEFAULT NULL,
  PRIMARY KEY (`perm_id`)
) ENGINE=InnoDB AUTO_INCREMENT=91 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permissions`
--

LOCK TABLES `permissions` WRITE;
/*!40000 ALTER TABLE `permissions` DISABLE KEYS */;
INSERT INTO `permissions` VALUES (1,1,1,1),(2,2,1,1),(3,3,1,1),(4,4,1,1),(5,5,1,1),(6,1,2,0),(7,2,2,7),(8,3,2,7),(9,4,2,7),(10,5,2,7),(11,1,3,0),(12,2,3,1),(13,3,3,1),(14,4,3,1),(15,5,3,1),(16,1,4,0),(17,2,4,1),(18,3,4,1),(19,4,4,1),(20,5,4,1),(21,1,5,0),(22,2,5,1),(23,3,5,1),(24,4,5,1),(25,5,5,1),(26,6,1,7),(27,7,1,7),(28,8,1,7),(29,9,1,7),(30,10,1,7),(31,11,1,7),(32,12,1,7),(33,13,1,7),(34,6,2,0),(35,7,2,1),(36,8,2,1),(37,9,2,1),(38,10,2,1),(39,11,2,1),(40,12,2,1),(41,13,2,1),(42,6,3,0),(43,7,3,1),(44,8,3,1),(45,9,3,1),(46,10,3,1),(47,11,3,0),(48,12,3,0),(49,13,3,0),(50,6,4,0),(51,7,4,1),(52,8,4,1),(53,9,4,3),(54,10,4,3),(55,11,4,1),(56,12,4,1),(57,13,4,1),(58,6,5,0),(59,7,5,1),(60,8,5,1),(61,9,5,0),(62,10,5,0),(63,11,5,0),(64,12,5,0),(65,13,5,0),(66,14,1,7),(67,15,1,7),(68,16,1,7),(69,17,1,7),(70,18,1,7),(71,14,2,0),(72,15,2,0),(73,16,2,0),(74,17,2,0),(75,18,2,0),(76,14,3,0),(77,15,3,0),(78,16,3,0),(79,17,3,0),(80,18,3,0),(81,14,4,0),(82,15,4,0),(83,16,4,2),(84,17,4,0),(85,18,4,0),(86,14,5,0),(87,15,5,0),(88,16,5,0),(89,17,5,0),(90,18,5,0);
/*!40000 ALTER TABLE `permissions` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-11-26 11:08:41
