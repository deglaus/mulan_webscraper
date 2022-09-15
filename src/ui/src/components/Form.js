import React, { Component } from "react";
import ScrapeSites from "../ExecuteScrapers.js";
// Routing manipulation support for this component
import { Redirect } from "react-router-dom";

// Used to send search string to other components no matter their hierarchy in app
import { subscriber } from "../messageService";
import "./Form.css";

class Form extends Component {
  constructor(props) {
    super(props);
    this.state = {
      // Initialize search string value
      searchString: "",
    };
  }
  handleSubmit = async (event) => {
    // Prevents default action going through
    event.preventDefault();
    try {
      var finalString = this.state.searchString;
      console.log(finalString);
      // Starts scrapers using API to communicate with the backend
      if (finalString.length > 0) {
        ScrapeSites(finalString);
        // Redirects to the result component using routes in App.js
        setTimeout(() => {
          this.setState({ redirect: "/results" });
        }, 15000);
        console.log("SEARCH STRING IN FORM:" + finalString);
        subscriber.next(finalString);
      }
    } catch (err) {
      console.log(err);
    }
  };

  handleInputChange = (event) => {
    event.preventDefault();
    // Update field value
    this.setState({
      [event.target.name]: event.target.value,
    });
  };

  render() {
    // Used for redirection when submitting search string
    if (this.state.redirect) {
      return <Redirect to={this.state.redirect} />;
    }
    return (
      <div>
        <form onSubmit={this.handleSubmit}>
          <p>
            <input
              className="FormInput"
              type="text"
              placeholder="What are you looking for?"
              name="searchString"
              onChange={this.handleInputChange}
            ></input>
          </p>
          <p>
            <button className="btn">Search</button>
          </p>
        </form>
      </div>
    );
  }
}

export default Form;
