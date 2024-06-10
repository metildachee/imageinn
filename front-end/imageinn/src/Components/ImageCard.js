import React, { useState, useEffect } from "react";
import "animate.css";
import { Col, Layout, Modal, Row, Space } from "antd";
import Title from "./Title";
import HorizontalScroll from "./HorizontalScroll";

const TILT_MIN_DEGREE = 5;
const TILT_MAX_DEGREE = 10;
const MARGIN = 5;
const BORDER_SIZE = 15;

const ImageCard = ({
  altText,
  imgSrc,
  imageCaption,
  categoryNames,
  score,
  id,
  userID,
}) => {
  const tiltDegree =
    Math.floor(Math.random() * (TILT_MAX_DEGREE - TILT_MIN_DEGREE + 1)) +
    TILT_MIN_DEGREE; // Random degree between TILT_MIN_DEGREE and TILT_MAX_DEGREE
  const boxShadow = `${tiltDegree / 2}px ${
    tiltDegree / 2
  }px 10px rgba(0, 0, 0, 0.5)`; // Casting shadow
  const borderStyle = `${BORDER_SIZE}px solid white`; // White border
  const colorPalette = ["#8A8FEA", "#DC648F", "#FFB6B9", "#FF8C61", "#FFC93C"]; // Define your color palette

  let color = colorPalette[Math.floor(Math.random() * colorPalette.length)];

  const [modelOpen, setModelOpen] = useState(false);
  const [images, setImages] = useState([]);

  const updateClick = async (e) => {
    e.stopPropagation(); // Stop the event from propagating
    console.log("Div clicked, opening modal");
    setModelOpen(true);

    // Call the API endpoint with the user_id and clicked image id
    const apiUrl = `http://localhost:5000/click/${userID}`;
    try {
      const response = await fetch(apiUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ clicked_items: [id] }),
      });
      const responseData = await response.json();
      console.log(`send click ${id, userID}: API response:`, responseData);
    } catch (error) {
      console.error("Error calling train_model API:", error);
    }
  };

  const callIDAPI = async (id) => {
    console.log("submit image input:", id);
    const url = `http://localhost:8080/search_by_id?id=${id}`;
    console.log(url);
    try {
      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      const responseData = await response.json();
      console.log("response from id API:", responseData);
      setImages(responseData.images.slice(1));
    } catch (error) {
      console.error("Error during API call:", error);
    }
  };

  useEffect(() => {
    callIDAPI(id);
  }, [id]);

  const decodeBase64 = (base64) => {
    try {
      const byteCharacters = atob(base64);
      const byteNumbers = new Array(byteCharacters.length);
      for (let i = 0; i < byteCharacters.length; i++) {
        byteNumbers[i] = byteCharacters.charCodeAt(i);
      }
      const byteArray = new Uint8Array(byteNumbers);
      const blob = new Blob([byteArray], { type: "image/jpeg" });
      return URL.createObjectURL(blob);
    } catch (e) {
      console.error(`Error decoding base64 string: ${e}`);
      return null;
    }
  };

  const handleCancel = (e) => {
    e.stopPropagation(); // Stop the event from propagating
    console.log("Modal close button clicked");
    setModelOpen(false);
  };

  const handleOpen = (e) => {
    e.stopPropagation(); // Stop the event from propagating
    console.log("Div clicked, opening modal");
    setModelOpen(true);
    updateClick(e);
  };

  return (
    <div
      style={{ transform: `rotate(${tiltDegree}deg)`, backgroundColor: "" }}
      onClick={handleOpen}
    >
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

      <Modal
        title=""
        centered
        open={modelOpen}
        onOk={handleCancel}
        onCancel={handleCancel} // Ensure onCancel closes the modal
        style={{
          // height: "60vh",
          borderRadius: "5px", // Adjust the border radius as needed
          border: "1px solid black",
        }}
        footer={null}
        width={1500}
        afterClose={(e) => e && e.stopPropagation()} // Stop the event from propagating after the modal is closed
      >
        <Layout style={{ backgroundColor: "#F7B796" }}>
          <Row>
            <Col span={8}>
              <div
                style={{
                  backgroundColor: "",
                  boxShadow,
                  border: borderStyle,
                  margin: MARGIN,
                  color: "#8A8FEA",
                  position: "relative", // Add position relative for positioning the paragraph
                  display: "block",
                }}
              >
                <img
                  src={imgSrc}
                  alt="current clicks"
                  style={{
                    width: "100%",
                    height: "100%",
                    display: "block",
                    objectFit: "cover",
                  }}
                />
              </div>
            </Col>
            <Col span={8}>
              <Space align="center" style={{ padding: "10" }}>
                <Row align={"center"} style={{ backgroundColor: "" }}>
                  <Title text={imageCaption} style={{ marginTop: "30px" }} />
                </Row>
              </Space>
            </Col>
            <Col span={8}>
              <Space align="center" style={{ padding: "10" }}>
                <Row align={"center"}>
                  <HorizontalScroll data={images} />
                </Row>
              </Space>
            </Col>
          </Row>
        </Layout>
      </Modal>
    </div>
  );
};

export default ImageCard;
