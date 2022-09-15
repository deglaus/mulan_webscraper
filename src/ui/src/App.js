import React from "react";

///////////////
/// ROUTING ///
///////////////
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import About from "./Routes/About";
import ScanCov from "./Routes/ScanCov";
import Tutorial from "./Routes/Tutorial";
import Results from "./Routes/Results";

///////////////
/// STYLING ///
///////////////
import "./App.css";
import logo from "./assets/logo-blend.png";

//////////////////
/// COMPONENTS ///
//////////////////
import NavBar from "./components/NavBar";
import Form from "./components/Form";
// Array used for the navigation bar (navigation choices)
import Menu from "./components/Menu";

/* The App component is the main component in React.js */

function App() {
  return (
    // Enables Routing for whole app div block
    <Router>
      <div className="App">
        <NavBar MenuItems={Menu}></NavBar>
        <Switch>
          <Route exact path="/">
            <img className="welcome-logo" src={logo} alt=""></img>
            <header className="App-header">
              <p>Welcome to Second-Hand Scanner!</p>
              {/* Pass the placeholder and imported data to search bar component
              <Search
                placeholder="Search for a product from database..."
                data={data}
              /> */}
            </header>
            <header className="App-header">
              {/* Frontend to backend communication using API according to user input*/}
              <Form></Form>
            </header>
          </Route>
          <Route exact path="/about">
            <About></About>
          </Route>
          <Route exact path="/stores">
            <ScanCov></ScanCov>
          </Route>
          <Route exact path="/tutorial">
            <Tutorial></Tutorial>
          </Route>
          <Route exact path="/results">
            <Results></Results>
          </Route>
        </Switch>
      </div>
    </Router>
  );
}

export default App;
