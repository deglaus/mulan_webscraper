import CollectFromTable from "./FetchFromDatabase";

const FetchFromAllTables = function(itemName) {
  let data = [];
  // data = CollectFromTable(itemName, "Tradera", data);
  data = CollectFromTable(itemName, "Bokbörsen", data);
  data = CollectFromTable(itemName, "Adlibris", data);
  data = CollectFromTable(itemName, "Ebay", data);

  return data;
};

export default FetchFromAllTables;
