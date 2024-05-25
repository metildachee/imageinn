import React, { useState, useEffect } from 'react';
import SearchBar from './SearchBar';
import ImageRow from './ImageRow';
import './App.css';
import CategoryButtonsRow from './CategoryButtonsRow';

function App() {
  const [images, setImages] = useState([]);  // State to store image data
  const [categories, setCategories] = useState([]);  // State to store image data

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('http://localhost:8080/category_info?count=5');
        const categories = await response.json();
        console.log("response from API:", categories);
        setCategories(categories.categories)
      } catch (error) {
        console.error('Error:', error);
      }
    };

    fetchData();
  }, []); // Empty dependency array ensures this effect runs once after the first render

  const onSearch = async (value) => {
    console.log("submit input:", value);
    try {
      const response = await fetch(`http://localhost:8080/search?q=${value}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        },
      });
      const responseData = await response.json();
      console.log("response from API:", responseData);
      setImages(responseData.images);  // Update state with the response data
    } catch (error) {
      console.error("Error during API call:", error);
    }
  };

  return (
    <div className="App">
        <SearchBar onSearch={onSearch} />
        {images.length > 0 ? <ImageRow images={images} /> : <div>No images to display</div>}
        {categories.length > 0 && <CategoryButtonsRow categories={categories} />}
    </div>
  );
}

export default App;
