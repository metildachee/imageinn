import React from "react";

const styles = {
  boldText: {
    fontFamily: "Roboto",
    fontWeight: "bold",
    fontSize: 40,
    color: "white",
  },
  italicText: {
    fontFamily: "Cormorant Garamond",
    fontStyle: "italic",
    fontWeight: "900",
    fontSize: 80,
    lineHeight: 1,
    color: "#F06126"
  },
};

const splitCaption = (caption) => {
  const index =
    caption.indexOf(",") !== -1 ? caption.indexOf(",") : caption.indexOf(" ");
  if (index !== -1) {
    const firstHalf = caption.substring(0, index);
    const secondHalf = caption.substring(index + 1);
    console.log(caption)
    return (
        <>
          <div style={{ textAlign: "left", marginLeft: "20px", marginTop: "100px" }}>
            <div style={{ margin: 0, padding: 0 }}><span style={styles.italicText}>{firstHalf.toUpperCase()}</span></div>
            <div style={{ marginTop: -15, padding: 0 }}><span style={styles.boldText}>{secondHalf.toUpperCase()}</span></div>
          </div>
        </>
      );
      
  }

  return <p style={styles.boldText}>{caption}</p>;
};

const Title = ({ text }) => {
  return splitCaption(text);
};

export default Title;
