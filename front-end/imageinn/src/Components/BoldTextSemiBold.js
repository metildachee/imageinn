import React from "react";

const BoldTextSemiBold = ({ bold, semiBold, fontSize }) => {
  return (
    <>
      <span
        style={{ marginLeft: "5px", fontSize: fontSize }}
        className="roboto-bold-italic"
      >
        {bold}
      </span>
      <span
        style={{ fontSize: fontSize }}
        className="cormorant-garamond-semibold-italic"
      >
        {" "}
        {semiBold}
      </span>
    </>
  );
};

export default BoldTextSemiBold;
