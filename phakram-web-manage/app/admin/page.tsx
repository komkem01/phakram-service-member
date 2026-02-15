export default function AdminDashboard() {
  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-800 mb-6">Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Card สรุปข้อมูล */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-500 text-sm">สมาชิกทั้งหมด</p>
              <p className="text-2xl font-bold text-gray-800 mt-2">1,234</p>
            </div>
            <div className="bg-blue-100 rounded-full p-3">
              <span className="text-2xl">👥</span>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-500 text-sm">คำสั่งซื้อวันนี้</p>
              <p className="text-2xl font-bold text-gray-800 mt-2">45</p>
            </div>
            <div className="bg-green-100 rounded-full p-3">
              <span className="text-2xl">🛒</span>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-500 text-sm">สินค้าทั้งหมด</p>
              <p className="text-2xl font-bold text-gray-800 mt-2">567</p>
            </div>
            <div className="bg-purple-100 rounded-full p-3">
              <span className="text-2xl">📦</span>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-500 text-sm">รายได้วันนี้</p>
              <p className="text-2xl font-bold text-gray-800 mt-2">฿12,345</p>
            </div>
            <div className="bg-yellow-100 rounded-full p-3">
              <span className="text-2xl">💰</span>
            </div>
          </div>
        </div>
      </div>

      {/* ตารางคำสั่งซื้อล่าสุด */}
      <div className="bg-white rounded-lg shadow mt-6 p-6">
        <h2 className="text-xl font-bold text-gray-800 mb-4">คำสั่งซื้อล่าสุด</h2>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200">
                <th className="text-left py-3 px-4 text-gray-600 font-medium">เลขที่คำสั่งซื้อ</th>
                <th className="text-left py-3 px-4 text-gray-600 font-medium">ลูกค้า</th>
                <th className="text-left py-3 px-4 text-gray-600 font-medium">ยอดรวม</th>
                <th className="text-left py-3 px-4 text-gray-600 font-medium">สถานะ</th>
                <th className="text-left py-3 px-4 text-gray-600 font-medium">วันที่</th>
              </tr>
            </thead>
            <tbody>
              <tr className="border-b border-gray-100 hover:bg-gray-50">
                <td className="py-3 px-4">#ORD-001</td>
                <td className="py-3 px-4">สมชาย ใจดี</td>
                <td className="py-3 px-4">฿1,250</td>
                <td className="py-3 px-4">
                  <span className="bg-yellow-100 text-yellow-800 px-3 py-1 rounded-full text-sm">
                    รอชำระเงิน
                  </span>
                </td>
                <td className="py-3 px-4">15/02/2026</td>
              </tr>
              <tr className="border-b border-gray-100 hover:bg-gray-50">
                <td className="py-3 px-4">#ORD-002</td>
                <td className="py-3 px-4">สมหญิง รักสวย</td>
                <td className="py-3 px-4">฿2,500</td>
                <td className="py-3 px-4">
                  <span className="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm">
                    กำลังจัดส่ง
                  </span>
                </td>
                <td className="py-3 px-4">15/02/2026</td>
              </tr>
              <tr className="hover:bg-gray-50">
                <td className="py-3 px-4">#ORD-003</td>
                <td className="py-3 px-4">วิชัย มั่นคง</td>
                <td className="py-3 px-4">฿850</td>
                <td className="py-3 px-4">
                  <span className="bg-green-100 text-green-800 px-3 py-1 rounded-full text-sm">
                    จัดส่งแล้ว
                  </span>
                </td>
                <td className="py-3 px-4">14/02/2026</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
