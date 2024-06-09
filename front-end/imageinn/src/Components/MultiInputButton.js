import React, { useState } from "react";
import { Button, Modal } from "antd";
import { Col, Row, Space, Flex, Layout } from "antd";
import { SearchOutlined, AlignRightOutlined } from "@ant-design/icons";
import BoldTextSemiBold from "./BoldTextSemiBold";

const MultiInputButton = ({
  onClick,
  leftBold,
  leftSemibold,
  rightBold,
  backgroundColor,
}) => {
  return (
    <>
      <Button
        type="default"
        onClick={onClick}
        style={{
          backgroundColor: backgroundColor,
          color: "black",
          border: "1px solid black",
          borderRadius: "5px",
          width: "70ch",
          marginTop: "10px",
        }}
      >
        <div
          style={{
            position: "absolute",
            bottom: -5,
            right: -5,
            width: `70ch`,
            height: `30px`,
            borderRight: "1px solid black",
            borderBottom: "1px solid black",
            borderTopColor: "#DC648F",
            borderLeftColor: "#DC648F",
            boxSizing: "border-box",
            overflow: "hidden", // Hide the top and left borders
            pointerEvents: "none",
            margin: "-2px",
            borderBottomRightRadius: "5px",
          }}
        />
        <Row>
          <Col span={12}>
            <div style={{ textAlign: "left" }}>
              <SearchOutlined />
              {/* <span
                style={{ marginLeft: "5px", fontSize: fontSize }}
                className="roboto-bold-italic"
              >
                {leftBold}
              </span>
              <span 
              style={{fontSize: fontSize}}
            className="cormorant-garamond-semibold-italic">
                {" "}
                {leftSemibold}
              </span> */}
              <BoldTextSemiBold
                bold={leftBold}
                semiBold={leftSemibold}
                fontWeight={"600"}
              />
            </div>
          </Col>
          <Col span={12}>
            <div style={{ textAlign: "right" }}>
              <span
                className="roboto-bold-italic"
                style={{ textTransform: "uppercase", marginRight: "5px" }}
              >
                "{rightBold}"
              </span>
              <AlignRightOutlined />
            </div>
          </Col>
        </Row>
      </Button>
    </>
  );
};

export default MultiInputButton;
