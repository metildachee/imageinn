import React, { useState, useEffect } from "react";
import { Col, Row, Space, Flex, Layout } from "antd";
import SearchBar from "./Components/SearchBar";
import ImageRow from "./Components/ImageRow";
import "./App.css";
import Model from "./Components/Model";
import Filters from "./Components/Filters";

const { Header, Footer, Content } = Layout;

const headerStyle = {
  textAlign: "center",
  color: "#fff",
  height: "15vh",
  paddingInline: 100,
  backgroundColor: "#F8DCE5",
};
const contentStyle = {
  textAlign: "center",
  height: "85vh",
  color: "#fff",
  backgroundColor: "#F8DCE5",
};
const footerStyle = {
  textAlign: "center",
  color: "#fff",
  backgroundColor: "#F8DCE5",
  height: "0vh",
};
const layoutStyle = {
  borderRadius: 8,
  overflow: "hidden",
};

function App() {
  const [images, setImages] = useState([]); // State to store image data
  const [query, setQuery] = useState("");
  const [isFuzzy, setIsFuzzy] = useState(false);
  const [isAnd, setIsAnd] = useState(false);
  const [excludes, setExcludes] = useState([]);
  const [searchResultCount, setSearchResultCount] = useState(0); // State to store the number of search results
  const [isImg, setIsImg] = useState(false)

  const callAPI = async (query, isFuzzy, isAnd, excludes) => {
    console.log("submit input:", query, isFuzzy, isAnd, excludes);
    const url = `http://localhost:8080/search?q=${query}&is_fuzzy=${isFuzzy}&excludes=${excludes}&is_and=${isAnd}`;
    console.log(url);
    try {
      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      const responseData = await response.json();
      console.log("response from API:", responseData);
      setImages(responseData.images); // Update state with the response data
      setSearchResultCount(responseData.images.length); // Update the search result count
    } catch (error) {
      console.error("Error during API call:", error);
    }
  };

  const callImgAPI = async (q) => {
    console.log("submit image input:", q);
    const url = `http://localhost:8080/search_in_image?q=${query}`;
    console.log(url);
    try {
      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      const responseData = await response.json();
      console.log("response from img API:", responseData);
      setImages(responseData.images); // Update state with the response data
      setSearchResultCount(responseData.images.length); // Update the search result count
    } catch (error) {
      console.error("Error during API call:", error);
    }
  };

  const onQuerySearch = (value) => {
    console.log("query input", value);
    setQuery(value);
  };


  useEffect(() => {
    console.log("is img", isImg)
    if (isImg) {
      callImgAPI(query)
      return 
    } 
    if (query !== "" || isFuzzy || isAnd || (excludes && excludes.length > 0)) {
      callAPI(query, isFuzzy, isAnd, excludes);
    }


  }, [query, isFuzzy, isAnd, excludes, isImg]);

  return (
    <Flex gap="middle" wrap>
      <Layout style={layoutStyle}>
        <Header style={headerStyle}>
          <Row>
            <Col span={10}></Col>
            <Col span={12}>
              <Space align="center">
                <SearchBar style={{ margin: "0px" }} onSearch={onQuerySearch} />
                <Filters updateFuzzy={setIsFuzzy} updateIsAnd={setIsAnd} updateExcludes={setExcludes} updateImage={setIsImg}/>
              </Space>
            </Col>
            <Col span={2}></Col>
          </Row>
          <Row>
            <Col span={8}></Col>
            <Col span={8}>
              
            </Col>
            <Col span={8}></Col>
          </Row>
        </Header>
        <Content style={contentStyle}>
          <Space align="center">
            {images != null && images.length > 0 ? (
              <ImageRow images={images} />
            ) : (
              <div>No images to display</div>
            )}
          </Space>
        </Content>
        <Footer style={footerStyle}>
        </Footer>
      </Layout>
    </Flex>
  );
}

export default App;
