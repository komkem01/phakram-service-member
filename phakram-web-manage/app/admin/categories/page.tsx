"use client";

import { useState, useEffect } from "react";
import Loading from "@/components/admin/Loading";

interface Category {
  id: string;
  name: string;
  description: string;
  product_count: number;
  created_at: string;
}

export default function CategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setTimeout(() => {
      setCategories([
        { id: "1", name: "หมวดหมู่ 1", description: "รายละเอียดหมวดหมู่ 1", product_count: 25, created_at: "2026-01-15" },
        { id: "2", name: "หมวดหมู่ 2", description: "รายละเอียดหมวดหมู่ 2", product_count: 18, created_at: "2026-01-20" },
        { id: "3", name: "หมวดหมู่ 3", description: "รายละเอียดหมวดหมู่ 3", product_count: 32, created_at: "2026-02-01" },
      ]);
      setIsLoading(false);
    }, 500);
  }, []);

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-800">จัดการหมวดหมู่</h1>
        <button className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors">
          + เพิ่มหมวดหมู่
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {categories.map((category) => (
          <div key={category.id} className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition-shadow">
            <div className="flex items-start justify-between mb-4">
              <div>
                <h3 className="text-xl font-bold text-gray-800">{category.name}</h3>
                <p className="text-sm text-gray-600 mt-1">{category.description}</p>
              </div>
            </div>
            <div className="flex items-center justify-between pt-4 border-t border-gray-200">
              <span className="text-sm text-gray-600">
                สินค้า: <span className="font-semibold text-blue-600">{category.product_count}</span> รายการ
              </span>
              <div className="flex gap-2">
                <button className="text-blue-600 hover:text-blue-700 p-2 rounded hover:bg-blue-50 transition-colors">
                  แก้ไข
                </button>
                <button className="text-red-600 hover:text-red-700 p-2 rounded hover:bg-red-50 transition-colors">
                  ลบ
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
