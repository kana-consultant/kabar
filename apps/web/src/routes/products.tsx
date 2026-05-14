import { createFileRoute } from '@tanstack/react-router';
import Products from "@/pages/home/Products";

export const Route = createFileRoute("/products")({
    component: Products,
});