import React from "react";
import "./Results.css";
// Fetch data from a singular table/store
import CollectFromTable from "../FetchFromDatabase";
import { subscriber } from "../messageService";

function Results() {
  console.log("LOADING COMPONENT...");
  // Variable for search string to receive from Form component
  let finalString = "";

  // Data from tables, initially empty arrays...
  let adlibrisData = [];
  let biblioData = [];
  let blocketData = [];
  let bokbörsenData = [];
  let citiboardData = [];
  let etsyData = [];
  let facebookData = [];
  let traderaData = [];

  // Array used to store all items sorted
  let sortedData = [];

  // Perform side effect from function component using React hook
  console.log("JUST LOADED RESULTS COMPONENT");
  subscriber.subscribe((v) => {
    finalString = "" + v;
  });

  // Fetch items for each store
  adlibrisData = CollectFromTable(finalString, "Adlibris", adlibrisData);
  biblioData = CollectFromTable(finalString, "Biblio", biblioData);
  blocketData = CollectFromTable(finalString, "Blocket", blocketData);
  bokbörsenData = CollectFromTable(finalString, "Bokbörsen", bokbörsenData);
  citiboardData = CollectFromTable(finalString, "Citiboard", citiboardData);
  etsyData = CollectFromTable(finalString, "Etsy", etsyData);
  facebookData = CollectFromTable(finalString, "FacebookMarket", facebookData);
  traderaData = CollectFromTable(finalString, "Tradera", traderaData);

  console.log("SEARCH STRING IN RESULTS:" + finalString);

  function sortData() {
    // While loop variable
    let i = 0;

    // Max number of items loaded is the while-condition
    while (i <= 30) {
      console.log("ENTERED WHILE-LOOP");
      if (adlibrisData.length !== 0) {
        sortedData.push(adlibrisData[i]);
      }
      if (biblioData.length !== 0) {
        sortedData.push(biblioData[i]);
      }
      if (blocketData.length !== 0) {
        sortedData.push(blocketData[i]);
      }
      if (bokbörsenData.length !== 0) {
        sortedData.push(bokbörsenData[i]);
      }
      if (citiboardData.length !== 0) {
        sortedData.push(citiboardData[i]);
      }
      if (etsyData.length !== 0) {
        sortedData.push(etsyData[i]);
      }
      if (facebookData.length !== 0) {
        sortedData.push(facebookData[i]);
      }
      if (traderaData.length !== 0) {
        sortedData.push(traderaData[i]);
      }
      i++;
    }
  }

  // Sort all the data collected from tables
  sortData();

  // REMOVE ALL UNDEFINED VALUES IN ARRAY
  sortedData = sortedData.filter(function(element) {
    return element !== undefined;
  });

  console.log(sortedData);
  return (
    <div className="Results">
      <div className="Results-content">
        <div className="content-container">
          <ul className="content">
            {/* RENDER THE ARRAY CONTAINING ALL ITEMS (SORTED) */}
            {sortedData.map((item) => {
              return (
                <li className="li-result">
                  {/* ---METADATA EXTRACTION--- */}
                  <a href={item.URL}>
                    <div className="productTitle">
                      {item.Title.substring(0, 50)}... <br></br>
                    </div>
                    {/* Get product picture */}
                    <img
                      src={item.pictureURL}
                      className="item-picture"
                      span="PICTURE MISSING"
                      alt=""
                    ></img>
                  </a>
                  <p className="source">Source: {item.Site}</p>
                  {/* '.substring(0,30)' Sets a 30char limit to the description */}
                  <p className="description">
                    Description: {item.Description.substring(0, 75)}
                  </p>
                  <p className="price">{item.Price.toFixed(2)}kr</p>
                </li>
              );
            })}
          </ul>
        </div>
      </div>
    </div>
  );
}

export default Results;
