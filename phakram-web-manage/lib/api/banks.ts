import { apiClient } from './client';

export interface Bank {
  id: string;
  name_th: string;
  name_abb_th: string;
  name_en: string;
  name_abb_en: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface BankListParams {
  page?: number;
  size?: number;
}

export interface BankCreateInput {
  name_th: string;
  name_abb_th: string;
  name_en: string;
  name_abb_en: string;
  is_active: boolean;
}

export interface BankUpdateInput {
  name_th: string;
  name_abb_th: string;
  name_en: string;
  name_abb_en: string;
  is_active: boolean;
}

interface BackendPaginate {
  page?: number;
  size?: number;
  total?: number;
  Page?: number;
  Size?: number;
  Total?: number;
}

export interface BankListResponse {
  data: Bank[];
  paginate: {
    page: number;
    size: number;
    total: number;
  };
}

export const banksApi = {
  async list(params?: BankListParams): Promise<BankListResponse> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.size) queryParams.append('size', params.size.toString());

    const query = queryParams.toString();
    const endpoint = `/system/banks${query ? `?${query}` : ''}`;

    const response = await apiClient.get<Bank[]>(endpoint);
    const rawPaginate = (response.paginate ?? {}) as BackendPaginate;

    return {
      data: response.data,
      paginate: {
        page: rawPaginate.page ?? rawPaginate.Page ?? 1,
        size: rawPaginate.size ?? rawPaginate.Size ?? params?.size ?? 10,
        total: rawPaginate.total ?? rawPaginate.Total ?? response.data.length,
      },
    };
  },

  async getById(id: string): Promise<Bank> {
    const response = await apiClient.get<Bank>(`/system/banks/${id}`);
    return response.data;
  },

  async create(data: BankCreateInput): Promise<void> {
    await apiClient.post('/system/banks', data);
  },

  async update(id: string, data: BankUpdateInput): Promise<void> {
    await apiClient.patch(`/system/banks/${id}`, data);
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/system/banks/${id}`);
  },
};
