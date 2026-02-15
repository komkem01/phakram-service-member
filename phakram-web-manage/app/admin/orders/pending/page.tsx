"use client";

import { useState, useEffect } from "react";
import Loading from "@/components/admin/Loading";

interface Order {
  id: string;
  order_number: string;
  customer: string;
  total: number;
  status: string;
  created_at: string;
}

export default function PendingOrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setTimeout(() => {
      setOrders([
        {
          id: "1",
          order_number: "ORD-001",
          customer: "สมชาย ใจดี",
          total: 1250,
          status: "pending",
          created_at: "2026-02-15 10:30",
        },
        {
          id: "2",
          order_number: "ORD-004",
          customer: "สมศรี รักงาม",
          total: 3500,
          status: "pending",
          created_at: "2026-02-15 11:45",
        },
      ]);
      setIsLoading(false);
    }, 500);
  }, []);

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-800">คำสั่งซื้อ - รอดำเนินการ</h1>
        <p className="text-gray-600 mt-2">รายการคำสั่งซื้อที่รอการดำเนินการ</p>
      </div>

      <div className="bg-white rounded-lg shadow">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 bg-gray-50">
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">เลขที่คำสั่งซื้อ</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">ลูกค้า</th>
                <th className="text-right py-4 px-6 text-gray-600 font-semibold">ยอดรวม</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">สถานะ</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">วันที่สั่งซื้อ</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">จัดการ</th>
              </tr>
            </thead>
            <tbody>
              {orders.map((order) => (
                <tr key={order.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="py-4 px-6 font-medium text-blue-600">#{order.order_number}</td>
                  <td className="py-4 px-6">{order.customer}</td>
                  <td className="py-4 px-6 text-right font-semibold">฿{order.total.toLocaleString()}</td>
                  <td className="py-4 px-6 text-center">
                    <span className="bg-orange-100 text-orange-800 px-3 py-1 rounded-full text-sm">
                      รอดำเนินการ
                    </span>
                  </td>
                  <td className="py-4 px-6 text-gray-600">{order.created_at}</td>
                  <td className="py-4 px-6">
                    <div className="flex items-center justify-center gap-2">
                      <button className="text-green-600 hover:text-green-700 px-3 py-1 rounded hover:bg-green-50 transition-colors font-medium">
                        ดำเนินการ
                      </button>
                      <button className="text-blue-600 hover:text-blue-700 px-3 py-1 rounded hover:bg-blue-50 transition-colors">
                        ดูรายละเอียด
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
