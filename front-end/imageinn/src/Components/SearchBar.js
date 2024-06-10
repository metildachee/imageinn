import React, { useState } from "react";
import { Input } from "antd";
import { UserOutlined, PlusOutlined, SearchOutlined, AlignRightOutlined } from "@ant-design/icons";

const SearchBar = ({ onSearch }) => {
  const [inputValue, setInputValue] = useState("");

  return (
    <div className="custom-input-container roboto-medium-italic" style={{ fontStyle: "Roboto"}}>
      <Input
        className="custom-input-field"
        placeholder="START SEARCHING!"
        // prefix={<UserOutlined className="custom-icon" />}
        prefix={<SearchOutlined style={{color: "black"}}/>}
        suffix={<span style={{textAlign: "right", color: "black"}}><AlignRightOutlined /></span>}
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        onKeyDown={(event) => {
          if (event.key === "Enter") {
            onSearch(inputValue);
          }
        }}
      />
    <div
          style={{
            position: "absolute",
            bottom: 6,
            right: -3,
            width: `25ch`,
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
    </div>
    
    
  );
};

export default SearchBar;

