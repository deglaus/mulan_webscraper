import React from "react";
import "./Tutorial.css";
import gif1 from "../assets/tutorial/step-one.gif";
import gif2 from "../assets/tutorial/step-two.gif";
import gif3 from "../assets/tutorial/step-three.gif";

function Tutorial() {
  return (
    <div className="App">
      <header className="App-header">
        <p>Tutorial</p>
        {/* STEP ONE */}
        <div className="separator">
          <p className="instruction">
            1.Type what you are looking for
          </p>
          <img src={gif1} className="gif"></img>
        </div>
        {/* STEP TWO */}
        <div className="separator">
          <p className="instruction">2. Choose a product</p>
          <img src={gif2} className="gif"></img>
        </div>
        {/* STEP THREE */}
        <div className="separator">
          <p className="instruction">
            3. Click on the title or the product picture in order to navigate to the
            source
          </p>
          <img src={gif3} className="gif"></img>
        </div>
      </header>
    </div>
  );
}

export default Tutorial;
