import React from "react";
import "./ScanCov.css";

////////////////////////
// IMPORT STORE LOGOS //
////////////////////////
import adlibris from "../assets/stores/adlibris.png";
import biblio from "../assets/stores/biblio.png";
import blocket from "../assets/stores/blocket.png";
import bokbörsen from "../assets/stores/bokbörsen.png";
import citiboard from "../assets/stores/citiboard.png";
import etsy from "../assets/stores/etsy.png";
import facebook from "../assets/stores/facebook.png";
import tradera from "../assets/stores/tradera.png";

function ScanCov() {
  return (
    <div className="content">
      <p className="Cov-header">Supported Second-Hand Stores</p>
      <div className="logos">
        {/* ADLIBRIS */}
        <a href="https://www.adlibris.com/se">
          <img
            src={adlibris}
            alt="adlibris-logo"
            className="square-store-logo"
          ></img>
        </a>
        {/* BIBLIO */}
        <a href="https://www.biblio.com/">
          <img
            src={biblio}
            alt="biblio-logo"
            className="rectangular-store-logo"
          ></img>
        </a>
        {/* BLOCKET */}
        <a href="https://www.blocket.se/">
          <img
            src={blocket}
            alt="blocket-logo"
            className="square-store-logo"
          ></img>
        </a>
        {/* BOKBÖRSEN */}
        <a href="https://www.bokborsen.se/">
          <img
            src={bokbörsen}
            alt="bokbörsen-logo"
            className="square-store-logo"
          ></img>
        </a>
        {/* CITIBOARD */}
        <a href="https://citiboard.se/">
          <img
            src={citiboard}
            alt="citiboard-logo"
            className="rectangular-store-logo"
          ></img>
        </a>
        {/* ETSY */}
        <a href="https://www.etsy.com/">
          <img src={etsy} alt="etsy-logo" className="square-store-logo"></img>
        </a>
        {/* FACEBOOK MARKETPLACE */}
        <a href="https://www.facebook.com/marketplace/">
          <img
            src={facebook}
            alt="facebook-logo"
            className="square-store-logo"
          ></img>
        </a>
        {/* TRADERA */}
        <a href="https://www.tradera.com/">
          <img
            src={tradera}
            alt="tradera-logo"
            className="rectangular-store-logo"
          ></img>
        </a>
      </div>
    </div>
  );
}

export default ScanCov;
