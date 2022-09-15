import React from "react";
import "./About.css";
import logo from "../assets/logo-blend.png";

function About() {
  return (
    <div className="About">
      <div className="About-content">
        <p className="About-header">Story:</p>
        <p className="About-paragraph">
          The Second-Hand Store was created as a project for the course
          Operating Systems and Process-Oriented Programming course at Uppsala
          University.
        </p>
        <p className="About-header">Made By:</p>
        <ul className="About-paragraph">
          <li>Axel Nilsson</li>
          <li>Casper Brandt</li>
          <li>Christian Ocklind</li>
          <li>Christoffer Björklund</li>
          <li>Johan Söderström</li>
          <li>Niklas Gotowiec</li>
          <li>Tofte Tjörneryd</li>
          <li>William Borg</li>
        </ul>
      </div>
      <img src={logo} alt="logo" className="About-logo"></img>
    </div>
  );
}

export default About;
