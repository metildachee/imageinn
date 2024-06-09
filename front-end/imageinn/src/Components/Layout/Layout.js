import React from 'react';
import "./layout.css";

const Layout = ({ children, contentMaxWidth }) => {
  const customCSSProperties = {
    "--content-max-width": contentMaxWidth,
  };
  return (
    <div className="wrapper" aria-label="Web site content">
      <main
        className="content"
        aria-label="Principal content of the web page."
        style={customCSSProperties}
      >
        {children}
      </main>
    </div>
  );
};

export default Layout;
