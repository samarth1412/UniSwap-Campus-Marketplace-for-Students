import type { Listing } from '../types/listing';

export const MOCK_LISTINGS: Listing[] = [
  {
    id: 1,
    title: 'CS101 Textbook - Clean Copy',
    description: 'Gently used CS101 textbook, no markings. Perfect for first-years.',
    price: 799.0,
    category: 'Books',
    imageUrl: 'https://via.placeholder.com/320x200?text=CS+Book',
    sellerName: 'Alice',
  },
  {
    id: 2,
    title: 'Used Laptop for Coding',
    description: '15-inch laptop, great for coding and browsing. Battery ~4 hours.',
    price: 25000.0,
    category: 'Electronics',
    imageUrl: 'https://via.placeholder.com/320x200?text=Laptop',
    sellerName: 'Rahul',
  },
  {
    id: 3,
    title: 'Hostel Study Table',
    description: 'Compact study table that fits in hostel rooms. Minor scratches.',
    price: 1500.0,
    category: 'Furniture',
    imageUrl: 'https://via.placeholder.com/320x200?text=Table',
    sellerName: 'Sneha',
  },
  {
    id: 4,
    title: 'Noise-Cancelling Headphones',
    description: 'Perfect for the library. Includes carry case.',
    price: 3500.0,
    category: 'Electronics',
    imageUrl: 'https://via.placeholder.com/320x200?text=Headphones',
    sellerName: 'Karan',
  },
];

