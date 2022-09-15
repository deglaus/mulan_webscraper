const ScrapeSites = function(itemName) {
  itemName = itemName.replaceAll(" ", "_");

  var url =
    "execute?" +
    new URLSearchParams({
      itemname: itemName,
    }).toString();

  console.log(url);

  fetch(url).then((res) => res.json());
};

export default ScrapeSites;
