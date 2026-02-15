"use client";

import { useState, useEffect } from "react";
import Loading from "@/components/admin/Loading";

export default function DeliveredOrdersPage() {
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setTimeout(() => setIsLoading(false), 500);
  }, []);

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-800">คำสั่งซื้อ - จัดส่งแล้ว</h1>
        <p className="text-gray-600 mt-2">รายการคำสั่งซื้อที่จัดส่งเรียบร้อยแล้ว</p>
      </div>
      <div className="bg-white rounded-lg shadow p-6">
        <p className="text-gray-500 text-center py-8">ไม่มีคำสั่งซื้อที่จัดส่งแล้ว</p>
      </div>
    </div>
  );
}
