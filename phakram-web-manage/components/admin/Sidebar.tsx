"use client";

import { useSidebar } from "@/contexts/SidebarContext";
import Link from "next/link";
import { usePathname } from "next/navigation";

interface MenuItem {
  id: string;
  label: string;
  href?: string;
  icon?: string;
  children?: MenuItem[];
}

const menuItems: MenuItem[] = [
  {
    id: "dashboard",
    label: "ภาพรวม",
    icon: "📊",
    href: "/admin",
  },
  {
    id: "system-settings",
    label: "ตั้งค่าระบบ",
    icon: "⚙️",
    children: [
      { id: "genders", label: "เพศ", href: "/admin/genders" },
      { id: "prefixes", label: "คำนำหน้า", href: "/admin/prefixes" },
      { id: "provinces", label: "จังหวัด", href: "/admin/provinces" },
      { id: "districts", label: "อำเภอ", href: "/admin/districts" },
      { id: "sub-districts", label: "ตำบล", href: "/admin/sub-districts" },
      { id: "zipcodes", label: "รหัสไปรษณีย์", href: "/admin/zipcodes" },
      { id: "tiers", label: "ระดับสมาชิก", href: "/admin/tiers" },
      { id: "statuses", label: "สถานะ", href: "/admin/statuses" },
      { id: "banks", label: "ธนาคาร", href: "/admin/banks" },
    ],
  },
  {
    id: "product-management",
    label: "จัดการสินค้า",
    icon: "📦",
    children: [
      { id: "products", label: "สินค้า", href: "/admin/products" },
      { id: "categories", label: "หมวดหมู่", href: "/admin/categories" },
    ],
  },
  {
    id: "order-management",
    label: "จัดการคำสั่งซื้อ",
    icon: "🛒",
    children: [
      {
        id: "orders",
        label: "รายการคำสั่งซื้อ",
        children: [
          { id: "pending", label: "รอดำเนินการ", href: "/admin/orders/pending" },
          { id: "awaiting-payment", label: "รอชำระเงิน", href: "/admin/orders/awaiting-payment" },
          { id: "awaiting-shipment", label: "รอจัดส่ง", href: "/admin/orders/awaiting-shipment" },
          { id: "shipping", label: "กำลังจัดส่ง", href: "/admin/orders/shipping" },
          { id: "delivered", label: "จัดส่งแล้ว", href: "/admin/orders/delivered" },
        ],
      },
      { id: "order-history", label: "ประวัติรายการคำสั่งซื้อ", href: "/admin/orders/history" },
    ],
  },
  {
    id: "user-management",
    label: "จัดการผู้ใช้งาน",
    icon: "👥",
    children: [
      { id: "members", label: "สมาชิก", href: "/admin/members" },
      { id: "admins", label: "ผู้ดูแลระบบ", href: "/admin/admins" },
    ],
  },
  {
    id: "admin-profile",
    label: "ตั้งค่าข้อมูลผู้ดูแลระบบ",
    icon: "👤",
    href: "/admin/profile",
  },
];

function MenuItemComponent({ item, level = 0 }: { item: MenuItem; level?: number }) {
  const { expandedMenus, toggleMenu, isCollapsed } = useSidebar();
  const pathname = usePathname();
  const isExpanded = expandedMenus.has(item.id);
  const isActive = item.href ? pathname === item.href : false;
  const hasChildren = item.children && item.children.length > 0;

  const paddingLeft = `${(level + 1) * 1}rem`;

  if (hasChildren) {
    return (
      <div>
        <button
          onClick={() => toggleMenu(item.id)}
          className={`w-full flex items-center justify-between px-4 py-3 text-left hover:bg-gray-100 transition-colors ${
            isActive ? "bg-blue-50 text-blue-600 font-medium" : "text-gray-700"
          }`}
          style={{ paddingLeft }}
        >
          <div className="flex items-center gap-3">
            {level === 0 && item.icon && <span className="text-xl">{item.icon}</span>}
            {!isCollapsed && <span className={level > 0 ? "text-sm" : ""}>{item.label}</span>}
          </div>
          {!isCollapsed && (
            <svg
              className={`w-4 h-4 transition-transform ${isExpanded ? "rotate-180" : ""}`}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
            </svg>
          )}
        </button>
        {isExpanded && !isCollapsed && (
          <div>
            {item.children?.map((child) => (
              <MenuItemComponent key={child.id} item={child} level={level + 1} />
            ))}
          </div>
        )}
      </div>
    );
  }

  return (
    <Link
      href={item.href || "#"}
      className={`flex items-center gap-3 px-4 py-3 hover:bg-gray-100 transition-colors ${
        isActive ? "bg-blue-50 text-blue-600 font-medium border-r-4 border-blue-600" : "text-gray-700"
      }`}
      style={{ paddingLeft }}
    >
      {level === 0 && item.icon && <span className="text-xl">{item.icon}</span>}
      {!isCollapsed && <span className={level > 0 ? "text-sm" : ""}>{item.label}</span>}
    </Link>
  );
}

export default function Sidebar() {
  const { isCollapsed, toggleSidebar } = useSidebar();

  return (
    <aside
      className={`bg-white border-r border-gray-200 h-screen sticky top-0 transition-all duration-300 ${
        isCollapsed ? "w-20" : "w-64"
      }`}
    >
      <div className="flex items-center justify-between p-4 border-b border-gray-200">
        {!isCollapsed && <h1 className="text-xl font-bold text-gray-800">Phakram Admin</h1>}
        <button
          onClick={toggleSidebar}
          className="p-2 rounded-lg hover:bg-gray-100 transition-colors"
          aria-label={isCollapsed ? "ขยาย sidebar" : "ย่อ sidebar"}
        >
          <svg
            className={`w-6 h-6 transition-transform ${isCollapsed ? "rotate-180" : ""}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M11 19l-7-7 7-7m8 14l-7-7 7-7"
            />
          </svg>
        </button>
      </div>

      <nav className="overflow-y-auto h-[calc(100vh-73px)]">
        {menuItems.map((item) => (
          <MenuItemComponent key={item.id} item={item} />
        ))}
      </nav>
    </aside>
  );
}
