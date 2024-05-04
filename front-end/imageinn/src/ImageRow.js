import React from 'react';
import ImageCard from './ImageCard'; // Make sure the path to ImageCard is correct

const ImageRow = ({ images }) => {
  return (
    <div style={{ display: 'flex', justifyContent: 'space-around', flexWrap: 'wrap' }}>
      {images.map((img, index) => (
        <ImageCard
          key={index}
          altText={img.caption}
          imgSrc={img.url}
          imageCaption={img.caption}
        />
      ))}
    </div>
  );
};

export default ImageRow;
