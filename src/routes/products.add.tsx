import { createFileRoute } from "@tanstack/react-router";
import { ProductForm } from "@/components/products/ProductForm/ProductForm";

export const Route = createFileRoute("/products/add")({
  component: AddProduct,
});

export default function AddProduct() {
  return (
    <>
      <ProductForm isEdit={false} />
    </>
  );
};