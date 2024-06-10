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
  height: "5vh",
  paddingInline: 100,
  backgroundColor: "#F8DCE5",
};
const contentStyle = {
  textAlign: "center",
  height: "75vh",
  color: "#fff",
  backgroundColor: "#F8DCE5",
};
const footerStyle = {
  textAlign: "center",
  color: "#fff",
  backgroundColor: "#F8DCE5",
  height: "20vh",
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
  const [isImg, setIsImg] = useState(false);
  const [userID, _] = useState(611);
  const [isPersonalise, setPersonalise] = useState(false);
  const [isNLP, setNLP] = useState(false);
  const [keywords, setKeywords] = useState([]);

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
      setSearchResultCount(responseData.total_count); // Update the search result count
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
      setSearchResultCount(responseData.total_count); // Update the search result count
    } catch (error) {
      console.error("Error during API call:", error);
    }
  };

  const batchGetAPI = async (ids) => {
    const url = `http://localhost:8080/batch_ids?ids=${ids}`;
    console.log(url);
    try {
      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      const responseData = await response.json();
      console.log("response from batch ids API:", responseData);
      setImages(responseData.images); // Update state with the response data
      setSearchResultCount(responseData.total_count); // Update the search result count
    } catch (error) {
      console.error("Error during API call:", error);
    }
  };

  const callPersonaliseAPI = async () => {
    const url = `http://localhost:5000/similar_items/${userID}`;
    console.log(url);
    try {
      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      const responseData = await response.json();
      console.log("response from personalised API:", responseData);
      // setImages(responseData.images); // Update state with the response data
      // setSearchResultCount(responseData.images.length); // Update the search result count
      return responseData.ids;
    } catch (error) {
      console.error("Error during API call:", error);
    }
  };

  const getNLPKeywords = async (query) => {
    const url = `http://localhost:5000/get_nlp_keywords/${query}`;
    console.log(url);
    try {
      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      const responseData = await response.json();
      console.log("response from nlp API:", responseData);
      // setImages(responseData.images); // Update state with the response data
      // setSearchResultCount(responseData.images.length); // Update the search result count
      setKeywords(responseData.keywords);
      return responseData.keywords;
    } catch (error) {
      console.error("Error during API call:", error);
    }
  };

  const onQuerySearch = (value) => {
    console.log("query input", value);
    setQuery(value);
  };

  useEffect(() => {
    console.log("is img", isImg);
    if (isImg) {
      callImgAPI(query);
      return;
    } else if (isPersonalise) {
      console.log("personalised query");
      callPersonaliseAPI().then((ids) => {
        console.log("batch get query");
        if (ids != null && ids.length > 0) {
          batchGetAPI(ids.slice(1, 5));
        }
      });
    } else if (isNLP && query) {
      console.log("nlp query");
      getNLPKeywords(query).then((keywords) => {
        callAPI(keywords, isFuzzy, isAnd, excludes);
      });
    } else if (
      query !== "" ||
      isFuzzy ||
      isAnd ||
      (excludes && excludes.length > 0)
    ) {
      callAPI(query, isFuzzy, isAnd, excludes);
    }
  }, [query, isFuzzy, isAnd, excludes, isImg, isPersonalise, isNLP]);

  return (
    <Flex gap="middle" wrap>
      <Layout style={layoutStyle}>
        <Header style={headerStyle}>
          <Row>
            <Col span={4}></Col>
            <Col span={16}>
              <div style={{ position: "relative" }}>
                <SearchBar style={{ margin: "0px" }} onSearch={onQuerySearch} />
                <span
                  style={{
                    position: "absolute",
                    bottom: -35,
                    left: 0,
                    color: "#8a8fea",
                    fontFamily: "Cormorant Garamond, serif",
                    fontSize: "18px",
                    margin: "0px",
                    padding: "5px",
                  }}
                >
                  {!(isNLP && keywords && keywords.length > 0)
                    ? `${searchResultCount} posters`
                    : `Searching on ${keywords}`}
                </span>
              </div>
            </Col>
            <Col span={4}></Col>
          </Row>
        </Header>
        <Content style={contentStyle}>
          <Space align="center">
            {images != null && images.length > 0 ? (
              <ImageRow userID={userID} images={images} />
            ) : (
              <div>No images to display</div>
            )}
          </Space>
        </Content>
        <Footer style={footerStyle}>
          <Filters
            updateFuzzy={setIsFuzzy}
            updateIsAnd={setIsAnd}
            updateExcludes={setExcludes}
            updateImage={setIsImg}
            updatePersonalise={setPersonalise}
            updateNLP={setNLP}
          />
        </Footer>
      </Layout>
    </Flex>
  );
}

export default App;
