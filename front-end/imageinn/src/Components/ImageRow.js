import React from "react";
import ImageCard from "./ImageCard";
import Layout from "./Layout/Layout";
import GalleryDiv from "./Layout/GalleryDiv";
import { Space } from "antd";

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

const ImageRow = ({ images }) => {
  return (
    <Layout>
        <GalleryDiv>
          {images.map((img, index) => (
            <ImageCard
              key={index}
              altText={img.caption}
              imgSrc={decodeBase64(img.img)}
              imageCaption={img.title}
              score={img.score}
              categoryNames={img.category_names}
            />
          ))}
        </GalleryDiv>
    </Layout>
  );
};

export default ImageRow;
