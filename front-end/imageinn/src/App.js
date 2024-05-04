import React, { useState } from 'react';
import SearchBar from './SearchBar';
import ImageRow from './ImageRow';
import './App.css';

function App() {
  const [images, setImages] = useState([]);  // State to store image data

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
        
    </div>
  );
}

export default App;
