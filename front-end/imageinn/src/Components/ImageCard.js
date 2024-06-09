import React from "react";
import "animate.css";

const TILT_MIN_DEGREE = 5;
const TILT_MAX_DEGREE = 10;
const MARGIN = 5;
const BORDER_SIZE = 15;

const ImageCard = ({ altText, imgSrc, imageCaption, categoryNames, score }) => {
  const tiltDegree =
    Math.floor(Math.random() * (TILT_MAX_DEGREE - TILT_MIN_DEGREE + 1)) +
    TILT_MIN_DEGREE; // Random degree between TILT_MIN_DEGREE and TILT_MAX_DEGREE
  const boxShadow = `${tiltDegree / 2}px ${
    tiltDegree / 2
  }px 10px rgba(0, 0, 0, 0.5)`; // Casting shadow
  const borderStyle = `${BORDER_SIZE}px solid white`; // White border
  const colorPalette = ["#8A8FEA", "#DC648F", "#FFB6B9", "#FF8C61", "#FFC93C"]; // Define your color palette

  let color = colorPalette[Math.floor(Math.random() * colorPalette.length)];

  return (
    <div style={{ transform: `rotate(${tiltDegree}deg)`, backgroundColor: "" }}>
      <div
        className="image-div on-hover"
        style={{
          backgroundColor: "",
          boxShadow,
          border: borderStyle,
          margin: MARGIN,
          color: "#8A8FEA",
          position: "relative", // Add position relative for positioning the paragraph
          // maxWidth: "300px", // Set width to 90%
        }}
      >
        <img
          src={imgSrc}
          alt="Image"
          style={{
            width: "100%",
            height: "auto",
            display: "block",
            objectFit: "cover",
          }}
        />
      </div>
      <span
        style={{
          // backgroundColor: ,
          fontFamily: "Roboto",
          fontWeight: 900,
          color: color,
          textAlign: "left",
          marginLeft: "25px",
          position: "absolute",
          bottom: 40, // Adjust as needed
          left: 0,
          right: 0,
        }}
      >
        {imageCaption.toUpperCase()}
      </span>
      <span
        style={{
          fontFamily: "Cormorant Garamond",
          fontWeight: 900,
          color: color,
          textAlign: "left",
          marginLeft: "25px",
          position: "absolute",
          bottom: 25, // Adjust as needed
          left: 0,
          right: 0,
          fontStyle: "italic",
          zIndex: 1,
        }}
      >
        {categoryNames.join(", ").toUpperCase()}
      </span>
      <span
        style={{
          position: "absolute",
          fontFamily: "cormorant-garamond-regular",
          top: 5,
          right: 5,
          backgroundColor: "rgba(0, 0, 0, 0.5)",
          color: "white",
          padding: "2px 5px",
          borderRadius: "3px",
        }}
      >
        Score: {score}
      </span>
    </div>
  );
};

export default ImageCard;
