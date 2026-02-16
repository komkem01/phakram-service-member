import { apiClient } from './client';

export interface Tier {
  id: string;
  name_th: string;
  name_en: string;
  min_spending: number;
  discount_rate: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface TierListParams {
  page?: number;
  size?: number;
}

export interface TierCreateInput {
  name_th: string;
  name_en: string;
  min_spending: string;
  discount_rate: string;
  is_active: boolean;
}

export interface TierUpdateInput {
  name_th: string;
  name_en: string;
  min_spending: string;
  discount_rate: string;
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

export interface TierListResponse {
  data: Tier[];
  paginate: {
    page: number;
    size: number;
    total: number;
  };
}

export const tiersApi = {
  async list(params?: TierListParams): Promise<TierListResponse> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.size) queryParams.append('size', params.size.toString());

    const query = queryParams.toString();
    const endpoint = `/system/tiers${query ? `?${query}` : ''}`;

    const response = await apiClient.get<Tier[]>(endpoint);
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

  async getById(id: string): Promise<Tier> {
    const response = await apiClient.get<Tier>(`/system/tiers/${id}`);
    return response.data;
  },

  async create(data: TierCreateInput): Promise<void> {
    await apiClient.post('/system/tiers', data);
  },

  async update(id: string, data: TierUpdateInput): Promise<void> {
    await apiClient.patch(`/system/tiers/${id}`, data);
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/system/tiers/${id}`);
  },
};
