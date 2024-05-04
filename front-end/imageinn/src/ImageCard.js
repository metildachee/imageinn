import React from 'react';
import { Card } from 'antd';
const { Meta } = Card;
const ImageCard = ({altText, imgSrc, imageCaption}) => (
  <Card
    hoverable
    style={{
      width: 240,
    }}
    cover={<img alt={altText} src={imgSrc} />}
  >
    <Meta title={imageCaption} description="www.instagram.com" />
  </Card>
);
export default ImageCard;