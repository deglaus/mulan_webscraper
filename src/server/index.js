const express = require("express");
const path = require("path");

const PORT = process.env.PORT || 3001;
const app = express();
var mysql = require("mysql");

var con = mysql.createConnection({
  host: "localhost",
  user: "admin",
  password: "mulan",
  database: "test",
});

con.connect(function (err) {
  if (err) throw err;
  console.log("Database is connected successfully !");
});

function getTable(data, callback) {
  let argument = "SELECT * FROM ".concat(data);
  con.query(argument, function (err, result, fields) {
    if (err) throw err;
    let resJSON = JSON.parse(JSON.stringify(result));
    return callback(resJSON);
  });
}

// Have Node serve the files for our built React app
app.use(express.static(path.resolve(__dirname, "../client/build")));

// Handle GET requests to /api route
app.get("/api", (req, res) => {
  var returnValue = "";
  let itemName = req.query.itemname;
  let storeName = req.query.storename;
  let tableName = storeName.concat("_" + itemName);
  console.log(tableName);
  getTable(tableName, function (result) {
    returnValue = result;

    res.json(returnValue);
  });
});

app.get("/execute", (req, res) => {
  let itemName = req.query.itemname;
  console.log("Scraping for " + itemName);

  const { exec } = require("child_process");
  exec(
    "cd .. && cd .. && make scrape arg=" + itemName,
    (err, stdout, stderr) => {
      if (err) {
        // node couldn't execute the command
        console.log(err);
        return;
      }

      // the *entire* stdout and stderr (buffered)
      console.log(`stdout: ${stdout}`);

      if (`${stderr}`.length != 0) {
        console.log(`stderr: ${stderr}`);
      }
      console.log("------------------------------------");
    }
  );
});

// All other GET requests not handled before will return our React app
app.get("*", (req, res) => {
  res.sendFile(path.resolve(__dirname, "../ui/public", "index.html"));
});

app.listen(PORT, () => {
  console.log(`Server listening on ${PORT}`);
});
