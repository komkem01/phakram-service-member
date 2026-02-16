"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import { Zipcode, zipcodesApi } from "@/lib/api/zipcodes";
import { subDistrictsApi } from "@/lib/api/sub-districts";
import { formatDateOnly } from "@/lib/utils/date";
import Loading from "@/components/admin/Loading";
import ZipcodeFormModal from "@/components/admin/ZipcodeFormModal";
import { useModal } from "@/contexts/ModalContext";

export default function ZipcodesPage() {
  const [zipcodes, setZipcodes] = useState<Zipcode[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [subDistrictNameMap, setSubDistrictNameMap] = useState<Record<string, string>>({});
  const subDistrictNameMapRef = useRef<Record<string, string>>({});
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState<'create' | 'edit'>('create');
  const [selectedZipcode, setSelectedZipcode] = useState<Zipcode | null>(null);

  const { showSuccess, showError, showWarning } = useModal();

  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  useEffect(() => {
    subDistrictNameMapRef.current = subDistrictNameMap;
  }, [subDistrictNameMap]);

  const hydrateMissingSubDistrictNames = useCallback(async (items: Zipcode[]) => {
    const neededIds = Array.from(
      new Set(
        items
          .map((item) => item.sub_districts_id ?? item.sub_district_id)
          .filter((value): value is string => Boolean(value))
      )
    );

    const missingIds = neededIds.filter((id) => !subDistrictNameMapRef.current[id]);
    if (missingIds.length === 0) return;

    const results = await Promise.allSettled(
      missingIds.map(async (id) => {
        const subDistrict = await subDistrictsApi.getById(id);
        return { id, name: subDistrict.name };
      })
    );

    const fallbackMap = results.reduce<Record<string, string>>((acc, result) => {
      if (result.status === 'fulfilled') {
        acc[result.value.id] = result.value.name;
      }
      return acc;
    }, {});

    if (Object.keys(fallbackMap).length > 0) {
      setSubDistrictNameMap((prev) => ({ ...prev, ...fallbackMap }));
    }
  }, []);

  const fetchSubDistricts = useCallback(async () => {
    try {
      const response = await subDistrictsApi.list({ page: 1, size: 1000 });
      const mapping = response.data.reduce<Record<string, string>>((acc, subDistrict) => {
        acc[subDistrict.id] = subDistrict.name;
        return acc;
      }, {});
      setSubDistrictNameMap(mapping);
    } catch (error) {
      console.error('Failed to load sub-districts:', error);
    }
  }, []);

  const fetchZipcodes = useCallback(async (currentPage: number, currentSize: number) => {
    try {
      setIsLoading(true);
      const response = await zipcodesApi.list({ page: currentPage, size: currentSize });
      setZipcodes(response.data);
      setTotal(response.paginate.total);
      setPage(response.paginate.page);
      await hydrateMissingSubDistrictNames(response.data);
    } catch {
      showError('ไม่สามารถโหลดข้อมูลรหัสไปรษณีย์ได้');
    } finally {
      setIsLoading(false);
    }
  }, [hydrateMissingSubDistrictNames, showError]);

  useEffect(() => {
    fetchZipcodes(page, pageSize);
  }, [fetchZipcodes, page, pageSize]);

  useEffect(() => {
    fetchSubDistricts();
  }, [fetchSubDistricts]);

  const handleCreate = () => {
    setModalMode('create');
    setSelectedZipcode(null);
    setIsModalOpen(true);
  };

  const handleEdit = (zipcode: Zipcode) => {
    setModalMode('edit');
    setSelectedZipcode(zipcode);
    setIsModalOpen(true);
  };

  const handleDelete = (zipcode: Zipcode) => {
    showWarning(
      `ต้องการลบรหัสไปรษณีย์ "${zipcode.name}" หรือไม่?`,
      async () => {
        try {
          await zipcodesApi.delete(zipcode.id);
          showSuccess('ลบรหัสไปรษณีย์สำเร็จ');
          fetchZipcodes(page, pageSize);
        } catch {
          showError('ไม่สามารถลบรหัสไปรษณีย์ได้');
        }
      },
      'ยืนยันการลบ'
    );
  };

  const handleSubmit = async (data: { sub_districts_id: string; name: string; is_active: boolean }) => {
    try {
      if (modalMode === 'create') {
        await zipcodesApi.create(data);
        showSuccess('เพิ่มรหัสไปรษณีย์สำเร็จ');
      } else if (selectedZipcode) {
        await zipcodesApi.update(selectedZipcode.id, data);
        showSuccess('แก้ไขรหัสไปรษณีย์สำเร็จ');
      }
      setIsModalOpen(false);
      fetchZipcodes(page, pageSize);
    } catch {
      showError(`ไม่สามารถ${modalMode === 'create' ? 'เพิ่ม' : 'แก้ไข'}รหัสไปรษณีย์ได้`);
    }
  };

  const handlePageChange = (newPage: number) => {
    if (newPage < 1 || newPage > totalPages) return;
    setPage(newPage);
  };

  const handlePageSizeChange = (newSize: number) => {
    setPageSize(newSize);
    setPage(1);
  };

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-800">จัดการรหัสไปรษณีย์</h1>
        <button
          onClick={handleCreate}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
        >
          + เพิ่มรหัสไปรษณีย์
        </button>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                รหัสไปรษณีย์
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                ตำบล
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                สถานะ
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                วันที่สร้าง
              </th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                จัดการ
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {zipcodes.map((zipcode) => (
              <tr key={zipcode.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {zipcode.name}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                  {subDistrictNameMap[zipcode.sub_districts_id ?? zipcode.sub_district_id ?? ''] ?? '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                      zipcode.is_active
                        ? 'bg-green-100 text-green-800'
                        : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {zipcode.is_active ? 'ใช้งาน' : 'ไม่ใช้งาน'}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatDateOnly(zipcode.created_at)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button
                    onClick={() => handleEdit(zipcode)}
                    className="text-blue-600 hover:text-blue-900 mr-4"
                  >
                    แก้ไข
                  </button>
                  <button
                    onClick={() => handleDelete(zipcode)}
                    className="text-red-600 hover:text-red-900"
                  >
                    ลบ
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between p-4 border-t border-gray-200">
          <div className="flex items-center gap-2 text-sm text-gray-600">
            <span>แสดง</span>
            <select
              value={pageSize}
              onChange={(e) => handlePageSizeChange(Number(e.target.value))}
              className="border border-gray-300 rounded-md px-2 py-1 bg-white"
            >
              <option value={10}>10</option>
              <option value={20}>20</option>
              <option value={50}>50</option>
              <option value={100}>100</option>
            </select>
            <span>รายการต่อหน้า</span>
          </div>

          <div className="text-sm text-gray-600">
            หน้า {page} / {totalPages} (ทั้งหมด {total} รายการ)
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={() => handlePageChange(page - 1)}
              disabled={page <= 1}
              className="px-3 py-1 rounded-md border border-gray-300 text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
            >
              ก่อนหน้า
            </button>
            <button
              onClick={() => handlePageChange(page + 1)}
              disabled={page >= totalPages}
              className="px-3 py-1 rounded-md border border-gray-300 text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
            >
              ถัดไป
            </button>
          </div>
        </div>
      </div>

      <ZipcodeFormModal
        isOpen={isModalOpen}
        mode={modalMode}
        zipcode={selectedZipcode}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
      />
    </div>
  );
}
