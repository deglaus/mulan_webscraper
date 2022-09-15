import React from "react";
import "./NavBar.css";

/* Navigation bar component */

/*Props */
// 1. Array including all the different choices for navigation.
// Such as the pages for Home, About us and Tutorial.
function NavBar({ MenuItems }) {
  return (
    <nav className="Selection">
      <ul className="Menu">
        {MenuItems.map((item, index) => {
          return (
            <li index={index}>
              <a className={item.className} href={item.url}>
                {item.title}
              </a>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}

export default NavBar;
