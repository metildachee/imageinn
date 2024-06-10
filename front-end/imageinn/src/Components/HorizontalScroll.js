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
        <Col span={2}>
          <div className="action-btns">
            <LeftCircleTwoTone twoToneColor="#22120A" onClick={() => handleScroll(-ITEM_WIDTH)} />
          </div>
        </Col>
        <Col span={20}>
          <div
            ref={containerRef}
            style={{
              width: "425px",
              overflowX: "scroll",
              scrollBehavior: "smooth",
            }}
          >
            <div className="content-box">
              {data && data.length > 0
                ? data.map((img) => (
                    <div className="card">
                      <img src={decodeBase64(img.img)} height="100%" />
                    </div>
                  ))
                : "no images to display"}
            </div>
          </div>
        </Col>
        <Col span={2}>
          <div className="action-btns">
            <RightCircleTwoTone twoToneColor="#22120A" onClick={() => handleScroll(ITEM_WIDTH)} />
          </div>
        </Col>
      </Row>
    </div>
  );
};

export default HorizontalScroll;
