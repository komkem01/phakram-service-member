import { apiClient } from './client';

export interface Prefix {
  id: string;
  name_th: string;
  name_en: string;
  gender_id: string;
  is_active: boolean;
  created_at: string;
}

export interface CreatePrefixDto {
  name_th: string;
  name_en: string;
  gender_id: string;
  is_active: boolean;
}

export interface UpdatePrefixDto {
  name_th: string;
  name_en: string;
  gender_id: string;
  is_active: boolean;
}

export interface ListPrefixesParams {
  page?: number;
  limit?: number;
}

export interface ListPrefixesResponse {
  data: Prefix[];
  paginate?: {
    page: number;
    size: number;
    total: number;
  };
}

export const prefixesApi = {
  /**
   * Get list of prefixes with pagination
   */
  async list(params?: ListPrefixesParams): Promise<ListPrefixesResponse> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());

    const query = queryParams.toString();
    const endpoint = `/system/prefixes/${query ? `?${query}` : ''}`;

    const response = await apiClient.get<Prefix[]>(endpoint);
    return {
      data: response.data,
      paginate: response.paginate,
    };
  },

  /**
   * Get single prefix by ID
   */
  async getById(id: string): Promise<Prefix> {
    const response = await apiClient.get<Prefix>(`/system/prefixes/${id}`);
    return response.data;
  },

  /**
   * Create new prefix
   */
  async create(data: CreatePrefixDto): Promise<void> {
    await apiClient.post('/system/prefixes/', data);
  },

  /**
   * Update existing prefix
   */
  async update(id: string, data: UpdatePrefixDto): Promise<void> {
    await apiClient.patch(`/system/prefixes/${id}`, data);
  },

  /**
   * Delete prefix
   */
  async delete(id: string): Promise<void> {
    await apiClient.delete(`/system/prefixes/${id}`);
  },
};
