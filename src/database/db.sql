--
-- Database: `db`
--
-- --------------------------------------------------------
--
-- Table structure for table `Blocket`
--
CREATE TABLE IF NOT EXISTS `Blocket` (
  `Title` varchar(30) NOT NULL COMMENT 'Put as primary for now, might change in the future',
  `Picture` varchar(2048) NOT NULL,
  `Description` text NOT NULL,
  `Price` int(11) NOT NULL,
  `SellerLocation` varchar(2048) NOT NULL,
  `DateAdded` varchar(30) NOT NULL,
  `Category` varchar(30) NOT NULL,
  `SellerAccount` varchar(2048) NOT NULL COMMENT 'Link to seller account ',
  `AdURL` varchar(2048) NOT NULL COMMENT 'Link to source '
) ENGINE = CSV DEFAULT CHARSET = utf8mb4 COMMENT = 'First draft of the Blocket data table. ';
-- --------------------------------------------------------
--
-- Table structure for table `Ebay`
--
CREATE TABLE IF NOT EXISTS `Ebay` (
  `Title` varchar(30) NOT NULL COMMENT 'Put as primary for now, might change in the future',
  `Picture` varchar(2048) NOT NULL,
  `Description` text NOT NULL,
  `Price` int(11) NOT NULL,
  `Quality` varchar(30) NOT NULL COMMENT 'The condition state of second hand product ',
  `SellerLocation` varchar(30) NOT NULL,
  `TimeLeft` int(11) NOT NULL,
  `Category` varchar(30) NOT NULL,
  `SellerAccount` varchar(2048) NOT NULL,
  `Views` int(11) NOT NULL,
  `Quantity` int(11) NOT NULL COMMENT 'Quantity still available ',
  `Returns` tinyint(1) NOT NULL,
  `Delivery` tinyint(1) NOT NULL,
  `AdURL` varchar(2048) NOT NULL
) ENGINE = CSV DEFAULT CHARSET = utf8mb4;
-- --------------------------------------------------------
--
-- Table structure for table `Facebook`
--
CREATE TABLE IF NOT EXISTS `Facebook` (
  `Title` varchar(30) NOT NULL,
  `Picture` varchar(2048) NOT NULL,
  `Description` text NOT NULL,
  `Price` int(11) NOT NULL,
  `SellerLocation` varchar(2048) NOT NULL,
  `FacebookUsername` varchar(100) NOT NULL,
  `Quality` varchar(30) NOT NULL,
  `SellerAccount` varchar(2048) NOT NULL,
  `AdURL` varchar(2048) NOT NULL COMMENT 'Link to the actual post, not the user account'
) ENGINE = CSV DEFAULT CHARSET = utf8mb4 COMMENT = 'Facebook Marketplace table ';
-- --------------------------------------------------------
--
-- Table structure for table `Tradera`
--
CREATE TABLE IF NOT EXISTS `Tradera` (
  `Title` varchar(30) NOT NULL,
  `Picture` varchar(2048) NOT NULL,
  `Description` text NOT NULL,
  `Price` int(11) NOT NULL,
  `SellerLocation` varchar(30) NOT NULL,
  `SellerAccount` varchar(2048) NOT NULL,
  `Delivery` tinyint(1) NOT NULL,
  `TimeLeft` int(11) NOT NULL,
  `DateAdded` varchar(30) NOT NULL,
  `Views` int(11) NOT NULL,
  `Category` varchar(30) NOT NULL,
  `Quality` varchar(30) NOT NULL,
  `AdURL` varchar(2048) NOT NULL
) ENGINE = CSV DEFAULT CHARSET = utf8mb4 COMMENT = 'Tradera Table ';
COMMIT;