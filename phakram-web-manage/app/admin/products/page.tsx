"use client";

import { useState, useEffect } from "react";
import Loading from "@/components/admin/Loading";

interface Product {
  id: string;
  name: string;
  price: number;
  stock: number;
  category: string;
  status: string;
}

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setTimeout(() => {
      setProducts([
        { id: "1", name: "สินค้า A", price: 1500, stock: 50, category: "หมวด 1", status: "active" },
        { id: "2", name: "สินค้า B", price: 2500, stock: 30, category: "หมวด 2", status: "active" },
        { id: "3", name: "สินค้า C", price: 800, stock: 0, category: "หมวด 1", status: "out_of_stock" },
      ]);
      setIsLoading(false);
    }, 500);
  }, []);

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-800">จัดการสินค้า</h1>
        <button className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors">
          + เพิ่มสินค้า
        </button>
      </div>

      <div className="bg-white rounded-lg shadow">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 bg-gray-50">
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">รหัส</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">ชื่อสินค้า</th>
                <th className="text-right py-4 px-6 text-gray-600 font-semibold">ราคา</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">คลังสินค้า</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">หมวดหมู่</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">สถานะ</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">จัดการ</th>
              </tr>
            </thead>
            <tbody>
              {products.map((product) => (
                <tr key={product.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="py-4 px-6 text-gray-600">#{product.id}</td>
                  <td className="py-4 px-6 font-medium">{product.name}</td>
                  <td className="py-4 px-6 text-right">฿{product.price.toLocaleString()}</td>
                  <td className="py-4 px-6 text-center">
                    <span className={product.stock === 0 ? "text-red-600 font-semibold" : ""}>
                      {product.stock}
                    </span>
                  </td>
                  <td className="py-4 px-6">{product.category}</td>
                  <td className="py-4 px-6 text-center">
                    <span
                      className={`px-3 py-1 rounded-full text-sm ${
                        product.status === "active"
                          ? "bg-green-100 text-green-800"
                          : "bg-red-100 text-red-800"
                      }`}
                    >
                      {product.status === "active" ? "พร้อมขาย" : "สินค้าหมด"}
                    </span>
                  </td>
                  <td className="py-4 px-6">
                    <div className="flex items-center justify-center gap-2">
                      <button className="text-blue-600 hover:text-blue-700 px-3 py-1 rounded hover:bg-blue-50 transition-colors">
                        แก้ไข
                      </button>
                      <button className="text-red-600 hover:text-red-700 px-3 py-1 rounded hover:bg-red-50 transition-colors">
                        ลบ
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
