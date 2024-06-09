// import React from "react"; // TODO: Change into a nicer category search bar
// import { AudioOutlined } from "@ant-design/icons";
// import { Input, Space } from "antd";
// const { Search } = Input;
// const suffix = (
//   <AudioOutlined
//     style={{
//       fontSize: 16,
//       color: "#353741",
//     }}
//   />
// );

// const SearchBar = ({ onSearch }) => (
//   // <Space direction="vertical">
//   //   <Search placeholder="input search text" onSearch={value => onSearch(value)} enterButton />
//   // </Space>

//   <div className="input-container">
//     <input
//       type="text"
//       className="input-field cormorant-garamond-regular"
//       placeholder="What's on your mind?"
//       // value={inputValue}
//       onKeyDown={(event) => {
//         if (event.key == "Enter") {
//           console.log(event)
//           onSearch(event.target.value);
//         }
//       }}
//     />
//   </div>
// );
// export default SearchBar;

import React, { useState } from "react";
import { Input } from "antd";
import { UserOutlined, PlusOutlined, SearchOutlined, AlignRightOutlined } from "@ant-design/icons";

const SearchBar = ({ onSearch }) => {
  const [inputValue, setInputValue] = useState("");

  return (
    <div className="custom-input-container roboto-medium-italic">
      <Input
        className="custom-input-field"
        placeholder="START SEARCHING!"
        // prefix={<UserOutlined className="custom-icon" />}
        prefix={<SearchOutlined />}
        suffix={<span style={{textAlign: "right"}}><AlignRightOutlined /></span>}
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        onKeyDown={(event) => {
          if (event.key === "Enter") {
            onSearch(inputValue);
            // setInputValue(""); // Clear the input field after searching
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

