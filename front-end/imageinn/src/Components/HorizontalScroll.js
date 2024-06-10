import React, { useRef, useState } from "react";
import { LeftCircleTwoTone, RightCircleTwoTone } from "@ant-design/icons";
import { Col, Row } from "antd";

const ITEM_WIDTH = 200;

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

const HorizontalScroll = ({ data }) => {
  const [scrollPosition, setScrollPosition] = useState(0);

  const containerRef = useRef();

  // Function to handle scrolling
  const handleScroll = (scrollAmount) => {
    const newScrollPosition = scrollPosition + scrollAmount;
    setScrollPosition(newScrollPosition);
    containerRef.current.scrollLeft = newScrollPosition;
  };

  return (
    <div className="scroll-container">
      <Row align="center">
        {/* <Col span={2}>
          <div className="action-btns" style={{justifyContent: "flex-start"}}>
            <LeftCircleTwoTone twoToneColor="#22120A" onClick={() => handleScroll(-ITEM_WIDTH)} />
          </div>
        </Col> */}
        <Col span={22}>
          <div
            ref={containerRef}
            style={{
              width: "200%",
              overflowX: "scroll",
              scrollBehavior: "smooth",
              marginLeft: "20px",
            }}
          >
            <div className="content-box">
              {data && data.length > 0
                ? data.map((img) => (
                    <div
                      className="card"
                      style={{ position: "relative", display: "inline-block" }}
                    >
                      <img
                        src={decodeBase64(img.img)}
                        height="100%"
                        alt="Image"
                      />
                      <span
                        style={{
                          position: "absolute",
                          bottom: 10,
                          left: 0,
                          padding: "5px",
                          fontFamily: "Roboto",
                          fontWeight: "900",
                          textTransform: "uppercase",
                        //   backgroundColor: "rgba(255, 255, 255, 0.7)",
                          color: "#E5D7D0"
                        }}
                      >
                        {img.title}
                      </span>
                      <span
                        style={{
                          position: "absolute",
                          bottom: -3,
                          left: -215,
                        //   transform: "translateX(-50%)",
                        //   padding: "5px",
                          fontFamily: "Garamond",
                        //   backgroundColor: "rgba(255, 255, 255, 0.7)",
                        }}
                      >
                        {img.category_names.length <= 2 ? img.category_names : img.category_names.slice(0, 2).join(", ")}
                      </span>
                    </div>
                  ))
                : "no images to display"}
            </div>
          </div>
        </Col>
        {/* <Col span={2}>
          <div className="action-btns" style= {{justifyContent: "flex-end"}} >
            <RightCircleTwoTone twoToneColor="#22120A" onClick={() => handleScroll(ITEM_WIDTH)} />
          </div>
        </Col> */}
      </Row>
    </div>
  );
};

export default HorizontalScroll;
