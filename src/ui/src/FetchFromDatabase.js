import React from "react";

const CollectFromTable = function(itemName, storeName, collectedData) {
  // Supports spaces in search string
  itemName = itemName.replaceAll(" ", "_");

  console.log("FETCH FROM DATABASE STRING:" + itemName);

  // Initialize data to an empty array
  const [data, setData] = React.useState([]);

  var url =
    "api?" +
    new URLSearchParams({
      itemname: itemName,
      storename: storeName,
    }).toString();

  console.log(url);

  React.useEffect(() => {
    fetch(url)
      .then((res) => res.json())
      .then((data) => setData(data));
  }, []);

  return collectedData.concat(data);
};

export default CollectFromTable;
