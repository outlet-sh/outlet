import webapi from './gocliRequest';
import * as components from './outletComponents';
export * from './outletComponents';

/**
 * @description
 * @param params
 */
export function adminListAgentActions(params: components.ListAgentActionsRequestParams) {
	return webapi.get<components.ListAgentActionsResponse>(`/api/admin/agent/actions`, params);
}

/**
 * @description
 * @param params
 */
export function adminApproveAgentAction(
	params: components.ApproveAgentActionRequestParams,
	id: string
) {
	return webapi.post<components.AgentActionLogEntry>(
		`/api/admin/agent/actions/${id}/approve`,
		params
	);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminRejectAgentAction(
	params: components.RejectAgentActionRequestParams,
	req: components.RejectAgentActionRequest,
	id: string
) {
	return webapi.post<components.AgentActionLogEntry>(
		`/api/admin/agent/actions/${id}/reject`,
		params,
		req
	);
}

/**
 * @description
 */
export function adminListCapabilities() {
	return webapi.get<components.ListCapabilitiesResponse>(`/api/admin/agent/capabilities`);
}

/**
 * @description
 * @param req
 */
export function adminUpdateCapabilities(req: components.UpdateCapabilitiesRequest) {
	return webapi.put<components.UpdateCapabilitiesResponse>(`/api/admin/agent/capabilities`, req);
}

/**
 * @description
 */
export function adminGetAgentConfig() {
	return webapi.get<components.AgentConfigInfo>(`/api/admin/agent/config`);
}

/**
 * @description
 * @param req
 */
export function adminUpdateAgentConfig(req: components.UpdateAgentConfigRequest) {
	return webapi.put<components.AgentConfigInfo>(`/api/admin/agent/config`, req);
}

/**
 * @description
 * @param params
 */
export function adminListAgentPrompts(params: components.ListAgentPromptsRequestParams) {
	return webapi.get<components.ListAgentPromptsResponse>(`/api/admin/agent/prompts`, params);
}

/**
 * @description
 * @param req
 */
export function adminCreateAgentPrompt(req: components.CreateAgentPromptRequest) {
	return webapi.post<components.AgentPromptInfo>(`/api/admin/agent/prompts`, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminUpdateAgentPrompt(
	params: components.UpdateAgentPromptRequestParams,
	req: components.UpdateAgentPromptRequest,
	id: string
) {
	return webapi.put<components.AgentPromptInfo>(`/api/admin/agent/prompts/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function adminDeleteAgentPrompt(
	params: components.DeleteAgentPromptRequestParams,
	id: string
) {
	return webapi.delete<components.EmptyResponse>(`/api/admin/agent/prompts/${id}`, params);
}

/**
 * @description
 */
export function listMCPAPIKeys() {
	return webapi.get<components.ListMCPAPIKeysResponse>(`/api/admin/api-keys`);
}

/**
 * @description
 * @param req
 */
export function createMCPAPIKey(req: components.CreateMCPAPIKeyRequest) {
	return webapi.post<components.CreateMCPAPIKeyResponse>(`/api/admin/api-keys`, req);
}

/**
 * @description
 * @param params
 */
export function revokeMCPAPIKey(params: components.RevokeMCPAPIKeyRequestParams, id: string) {
	return webapi.delete<components.RevokeMCPAPIKeyResponse>(`/api/admin/api-keys/${id}`, params);
}

/**
 * @description
 */
export function logout() {
	return webapi.post<components.AnalyticsResponse>(`/api/admin/auth/logout`);
}

/**
 * @description
 */
export function getCurrentUser() {
	return webapi.get<components.UserInfo>(`/api/admin/auth/me`);
}

/**
 * @description
 * @param params
 */
export function listProductBundles(params: components.ListProductBundlesRequestParams) {
	return webapi.get<components.ListProductBundlesResponse>(`/api/admin/billing/bundles`, params);
}

/**
 * @description
 * @param req
 */
export function createProductBundle(req: components.CreateProductBundleRequest) {
	return webapi.post<components.ProductBundleInfo>(`/api/admin/billing/bundles`, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateProductBundle(
	params: components.UpdateProductBundleRequestParams,
	req: components.UpdateProductBundleRequest,
	id: string
) {
	return webapi.put<components.ProductBundleInfo>(`/api/admin/billing/bundles/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteProductBundle(
	params: components.DeleteProductBundleRequestParams,
	id: string
) {
	return webapi.delete<null>(`/api/admin/billing/bundles/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function listProductCustomFields(params: components.ListProductCustomFieldsRequestParams) {
	return webapi.get<components.ListProductCustomFieldsResponse>(
		`/api/admin/billing/custom-fields`,
		params
	);
}

/**
 * @description
 * @param req
 */
export function createProductCustomField(req: components.CreateProductCustomFieldRequest) {
	return webapi.post<components.ProductCustomFieldInfo>(`/api/admin/billing/custom-fields`, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateProductCustomField(
	params: components.UpdateProductCustomFieldRequestParams,
	req: components.UpdateProductCustomFieldRequest,
	id: string
) {
	return webapi.put<components.ProductCustomFieldInfo>(
		`/api/admin/billing/custom-fields/${id}`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 */
export function deleteProductCustomField(
	params: components.DeleteProductCustomFieldRequestParams,
	id: string
) {
	return webapi.delete<null>(`/api/admin/billing/custom-fields/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function listPriceCurrencies(params: components.ListPriceCurrenciesRequestParams) {
	return webapi.get<components.ListPriceCurrenciesResponse>(
		`/api/admin/billing/price-currencies`,
		params
	);
}

/**
 * @description
 * @param req
 */
export function createPriceCurrency(req: components.CreatePriceCurrencyRequest) {
	return webapi.post<components.PriceCurrencyInfo>(`/api/admin/billing/price-currencies`, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updatePriceCurrency(
	params: components.UpdatePriceCurrencyRequestParams,
	req: components.UpdatePriceCurrencyRequest,
	id: string
) {
	return webapi.put<components.PriceCurrencyInfo>(
		`/api/admin/billing/price-currencies/${id}`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 */
export function deletePriceCurrency(
	params: components.DeletePriceCurrencyRequestParams,
	id: string
) {
	return webapi.delete<null>(`/api/admin/billing/price-currencies/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function listPriceTiers(params: components.ListPriceTiersRequestParams) {
	return webapi.get<components.ListPriceTiersResponse>(`/api/admin/billing/price-tiers`, params);
}

/**
 * @description
 * @param req
 */
export function createPriceTier(req: components.CreatePriceTierRequest) {
	return webapi.post<components.PriceTierInfo>(`/api/admin/billing/price-tiers`, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updatePriceTier(
	params: components.UpdatePriceTierRequestParams,
	req: components.UpdatePriceTierRequest,
	id: string
) {
	return webapi.put<components.PriceTierInfo>(`/api/admin/billing/price-tiers/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deletePriceTier(params: components.DeletePriceTierRequestParams, id: string) {
	return webapi.delete<null>(`/api/admin/billing/price-tiers/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function listBillingPrices(params: components.ListBillingPricesRequestParams) {
	return webapi.get<components.ListBillingPricesResponse>(`/api/admin/billing/prices`, params);
}

/**
 * @description
 * @param req
 */
export function createBillingPrice(req: components.CreateBillingPriceRequest) {
	return webapi.post<components.BillingPriceInfo>(`/api/admin/billing/prices`, req);
}

/**
 * @description
 * @param params
 */
export function getBillingPrice(params: components.GetBillingPriceRequestParams, id: string) {
	return webapi.get<components.BillingPriceInfo>(`/api/admin/billing/prices/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateBillingPrice(
	params: components.UpdateBillingPriceRequestParams,
	req: components.UpdateBillingPriceRequest,
	id: string
) {
	return webapi.put<components.BillingPriceInfo>(`/api/admin/billing/prices/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteBillingPrice(params: components.DeleteBillingPriceRequestParams, id: string) {
	return webapi.delete<null>(`/api/admin/billing/prices/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function listBillingProducts(params: components.ListBillingProductsRequestParams) {
	return webapi.get<components.ListBillingProductsResponse>(`/api/admin/billing/products`, params);
}

/**
 * @description
 * @param req
 */
export function createBillingProduct(req: components.CreateBillingProductRequest) {
	return webapi.post<components.BillingProductInfo>(`/api/admin/billing/products`, req);
}

/**
 * @description
 * @param params
 */
export function getBillingProduct(params: components.GetBillingProductRequestParams, id: string) {
	return webapi.get<components.BillingProductInfo>(`/api/admin/billing/products/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateBillingProduct(
	params: components.UpdateBillingProductRequestParams,
	req: components.UpdateBillingProductRequest,
	id: string
) {
	return webapi.put<components.BillingProductInfo>(
		`/api/admin/billing/products/${id}`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 */
export function deleteBillingProduct(
	params: components.DeleteBillingProductRequestParams,
	id: string
) {
	return webapi.delete<null>(`/api/admin/billing/products/${id}`, params);
}

/**
 * @description
 * @param req
 */
export function createQuoteItem(req: components.CreateQuoteItemRequest) {
	return webapi.post<components.CustomQuoteItemInfo>(`/api/admin/billing/quote-items`, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateQuoteItem(
	params: components.UpdateQuoteItemRequestParams,
	req: components.UpdateQuoteItemRequest,
	id: string
) {
	return webapi.put<components.CustomQuoteItemInfo>(
		`/api/admin/billing/quote-items/${id}`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 */
export function deleteQuoteItem(params: components.DeleteQuoteItemRequestParams, id: string) {
	return webapi.delete<null>(`/api/admin/billing/quote-items/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function listCustomQuotes(params: components.ListCustomQuotesRequestParams) {
	return webapi.get<components.ListCustomQuotesResponse>(`/api/admin/billing/quotes`, params);
}

/**
 * @description
 * @param req
 */
export function createCustomQuote(req: components.CreateCustomQuoteRequest) {
	return webapi.post<components.CustomQuoteInfo>(`/api/admin/billing/quotes`, req);
}

/**
 * @description
 * @param params
 */
export function getCustomQuote(params: components.GetCustomQuoteRequestParams, id: string) {
	return webapi.get<components.CustomQuoteInfo>(`/api/admin/billing/quotes/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateCustomQuote(
	params: components.UpdateCustomQuoteRequestParams,
	req: components.UpdateCustomQuoteRequest,
	id: string
) {
	return webapi.put<components.CustomQuoteInfo>(`/api/admin/billing/quotes/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteCustomQuote(params: components.DeleteCustomQuoteRequestParams, id: string) {
	return webapi.delete<null>(`/api/admin/billing/quotes/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function acceptCustomQuote(params: components.AcceptCustomQuoteRequestParams, id: string) {
	return webapi.post<components.CustomQuoteInfo>(`/api/admin/billing/quotes/${id}/accept`, params);
}

/**
 * @description
 */
export function syncBillingToStripe() {
	return webapi.post<components.SyncBillingResponse>(`/api/admin/billing/sync`);
}

/**
 * @description
 * @param params
 */
export function listCampaigns(params: components.ListCampaignsRequestParams) {
	return webapi.get<components.ListCampaignsResponse>(`/api/admin/campaigns`, params);
}

/**
 * @description
 * @param req
 */
export function createCampaign(req: components.CreateCampaignRequest) {
	return webapi.post<components.CampaignInfo>(`/api/admin/campaigns`, req);
}

/**
 * @description
 * @param params
 */
export function getCampaign(params: components.GetCampaignRequestParams, id: string) {
	return webapi.get<components.CampaignInfo>(`/api/admin/campaigns/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateCampaign(
	params: components.UpdateCampaignRequestParams,
	req: components.UpdateCampaignRequest,
	id: string
) {
	return webapi.put<components.CampaignInfo>(`/api/admin/campaigns/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteCampaign(params: components.DeleteCampaignRequestParams, id: string) {
	return webapi.delete<components.AnalyticsResponse>(`/api/admin/campaigns/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function scheduleCampaign(
	params: components.ScheduleCampaignRequestParams,
	req: components.ScheduleCampaignRequest,
	id: string
) {
	return webapi.post<components.CampaignInfo>(`/api/admin/campaigns/${id}/schedule`, params, req);
}

/**
 * @description
 * @param params
 */
export function sendCampaignNow(params: components.SendCampaignNowRequestParams, id: string) {
	return webapi.post<components.CampaignInfo>(`/api/admin/campaigns/${id}/send`, params);
}

/**
 * @description
 * @param params
 */
export function getCampaignStats(params: components.GetCampaignRequestParams, id: string) {
	return webapi.get<components.CampaignStatsResponse>(`/api/admin/campaigns/${id}/stats`, params);
}

/**
 * @description
 */
export function adminSiteAnalytics() {
	return webapi.get<components.AdminAnalyticsResponse>(`/api/admin/analytics`);
}

/**
 * @description
 */
export function adminListContacts() {
	return webapi.get<components.AdminContactsResponse>(`/api/admin/contacts`);
}

/**
 * @description
 * @param params
 */
export function adminMarkContactRead(params: components.MarkContactReadRequestParams, id: number) {
	return webapi.post<components.MarkContactReadResponse>(`/api/admin/contacts/${id}/read`, params);
}

/**
 * @description
 */
export function adminStats() {
	return webapi.get<components.AdminStatsResponse>(`/api/admin/stats`);
}

/**
 * @description
 * @param req
 */
export function createAuthorProfile(req: components.CreateAuthorProfileRequest) {
	return webapi.post<components.AuthorProfileInfo>(`/api/admin/content/authors`, req);
}

/**
 * @description
 * @param params
 */
export function listAuthorProfiles(params: components.ListAuthorProfilesRequestParams) {
	return webapi.get<components.ListAuthorProfilesResponse>(`/api/admin/content/authors`, params);
}

/**
 * @description
 * @param params
 */
export function getAuthorProfile(params: components.GetAuthorProfileRequestParams, id: string) {
	return webapi.get<components.AuthorProfileInfo>(`/api/admin/content/authors/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateAuthorProfile(
	params: components.UpdateAuthorProfileRequestParams,
	req: components.UpdateAuthorProfileRequest,
	id: string
) {
	return webapi.put<components.AuthorProfileInfo>(`/api/admin/content/authors/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteAuthorProfile(
	params: components.DeleteAuthorProfileRequestParams,
	id: string
) {
	return webapi.delete<null>(`/api/admin/content/authors/${id}`, params);
}

/**
 * @description
 * @param req
 */
export function createContentCategory(req: components.CreateContentCategoryRequest) {
	return webapi.post<components.ContentCategoryInfo>(`/api/admin/content/categories`, req);
}

/**
 * @description
 */
export function listContentCategories() {
	return webapi.get<components.ListContentCategoriesResponse>(`/api/admin/content/categories`);
}

/**
 * @description
 * @param params
 */
export function getContentCategory(params: components.GetContentCategoryRequestParams, id: string) {
	return webapi.get<components.ContentCategoryInfo>(`/api/admin/content/categories/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateContentCategory(
	params: components.UpdateContentCategoryRequestParams,
	req: components.UpdateContentCategoryRequest,
	id: string
) {
	return webapi.put<components.ContentCategoryInfo>(
		`/api/admin/content/categories/${id}`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 */
export function deleteContentCategory(
	params: components.DeleteContentCategoryRequestParams,
	id: string
) {
	return webapi.delete<null>(`/api/admin/content/categories/${id}`, params);
}

/**
 * @description
 * @param req
 */
export function createNavigationItem(req: components.CreateNavigationItemRequest) {
	return webapi.post<components.NavigationItemInfo>(`/api/admin/content/menu-items`, req);
}

/**
 * @description
 * @param params
 */
export function listNavigationItems(params: components.ListNavigationItemsRequestParams) {
	return webapi.get<components.ListNavigationItemsResponse>(
		`/api/admin/content/menu-items`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function getNavigationItem(params: components.GetNavigationItemRequestParams, id: string) {
	return webapi.get<components.NavigationItemInfo>(`/api/admin/content/menu-items/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateNavigationItem(
	params: components.UpdateNavigationItemRequestParams,
	req: components.UpdateNavigationItemRequest,
	id: string
) {
	return webapi.put<components.NavigationItemInfo>(
		`/api/admin/content/menu-items/${id}`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 */
export function deleteNavigationItem(
	params: components.DeleteNavigationItemRequestParams,
	id: string
) {
	return webapi.delete<null>(`/api/admin/content/menu-items/${id}`, params);
}

/**
 * @description
 * @param req
 */
export function createNavigationMenu(req: components.CreateNavigationMenuRequest) {
	return webapi.post<components.NavigationMenuInfo>(`/api/admin/content/menus`, req);
}

/**
 * @description
 * @param params
 */
export function listNavigationMenus(params: components.ListNavigationMenusRequestParams) {
	return webapi.get<components.ListNavigationMenusResponse>(`/api/admin/content/menus`, params);
}

/**
 * @description
 * @param params
 */
export function getNavigationMenu(params: components.GetNavigationMenuRequestParams, id: string) {
	return webapi.get<components.NavigationMenuInfo>(`/api/admin/content/menus/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateNavigationMenu(
	params: components.UpdateNavigationMenuRequestParams,
	req: components.UpdateNavigationMenuRequest,
	id: string
) {
	return webapi.put<components.NavigationMenuInfo>(`/api/admin/content/menus/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteNavigationMenu(
	params: components.DeleteNavigationMenuRequestParams,
	id: string
) {
	return webapi.delete<null>(`/api/admin/content/menus/${id}`, params);
}

/**
 * @description
 * @param req
 */
export function createContentPage(req: components.CreateContentPageRequest) {
	return webapi.post<components.ContentPageInfo>(`/api/admin/content/pages`, req);
}

/**
 * @description
 * @param params
 */
export function listContentPages(params: components.ListContentPagesRequestParams) {
	return webapi.get<components.ListContentPagesResponse>(`/api/admin/content/pages`, params);
}

/**
 * @description
 * @param params
 */
export function getContentPage(params: components.GetContentPageRequestParams, id: string) {
	return webapi.get<components.ContentPageInfo>(`/api/admin/content/pages/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateContentPage(
	params: components.UpdateContentPageRequestParams,
	req: components.UpdateContentPageRequest,
	id: string
) {
	return webapi.put<components.ContentPageInfo>(`/api/admin/content/pages/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteContentPage(params: components.DeleteContentPageRequestParams, id: string) {
	return webapi.delete<null>(`/api/admin/content/pages/${id}`, params);
}

/**
 * @description
 * @param req
 */
export function createContentPost(req: components.CreateContentPostRequest) {
	return webapi.post<components.ContentPostInfo>(`/api/admin/content/posts`, req);
}

/**
 * @description
 * @param params
 */
export function listContentPosts(params: components.ListContentPostsRequestParams) {
	return webapi.get<components.ListContentPostsResponse>(`/api/admin/content/posts`, params);
}

/**
 * @description
 * @param params
 */
export function getContentPost(params: components.GetContentPostRequestParams, id: string) {
	return webapi.get<components.ContentPostInfo>(`/api/admin/content/posts/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateContentPost(
	params: components.UpdateContentPostRequestParams,
	req: components.UpdateContentPostRequest,
	id: string
) {
	return webapi.put<components.ContentPostInfo>(`/api/admin/content/posts/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteContentPost(params: components.DeleteContentPostRequestParams, id: string) {
	return webapi.delete<null>(`/api/admin/content/posts/${id}`, params);
}

/**
 * @description
 * @param req
 */
export function createPageTemplate(req: components.CreatePageTemplateRequest) {
	return webapi.post<components.PageTemplateInfo>(`/api/admin/content/templates`, req);
}

/**
 * @description
 */
export function listPageTemplates() {
	return webapi.get<components.ListPageTemplatesResponse>(`/api/admin/content/templates`);
}

/**
 * @description
 * @param params
 */
export function getPageTemplate(params: components.GetPageTemplateRequestParams, id: string) {
	return webapi.get<components.PageTemplateInfo>(`/api/admin/content/templates/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updatePageTemplate(
	params: components.UpdatePageTemplateRequestParams,
	req: components.UpdatePageTemplateRequest,
	id: string
) {
	return webapi.put<components.PageTemplateInfo>(`/api/admin/content/templates/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deletePageTemplate(params: components.DeletePageTemplateRequestParams, id: string) {
	return webapi.delete<null>(`/api/admin/content/templates/${id}`, params);
}

/**
 * @description
 */
export function adminListCustomerTags() {
	return webapi.get<components.ListCustomerTagsResponse>(`/api/admin/customer-tags`);
}

/**
 * @description
 * @param req
 */
export function adminCreateCustomerTag(req: components.CreateCustomerTagRequest) {
	return webapi.post<components.CustomerTagInfo>(`/api/admin/customer-tags`, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminUpdateCustomerTag(
	params: components.UpdateCustomerTagRequestParams,
	req: components.UpdateCustomerTagRequest,
	id: string
) {
	return webapi.put<components.CustomerTagInfo>(`/api/admin/customer-tags/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function adminDeleteCustomerTag(
	params: components.DeleteCustomerTagRequestParams,
	id: string
) {
	return webapi.delete<components.TagOperationResponse>(`/api/admin/customer-tags/${id}`, params);
}

/**
 * @description
 * @param req
 */
export function adminAssignTagToCustomer(req: components.AssignTagToCustomerRequest) {
	return webapi.post<components.TagOperationResponse>(`/api/admin/customer-tags/assign`, req);
}

/**
 * @description
 * @param req
 */
export function adminRemoveTagFromCustomer(req: components.RemoveTagFromCustomerRequest) {
	return webapi.post<components.TagOperationResponse>(`/api/admin/customer-tags/remove`, req);
}

/**
 * @description
 * @param params
 */
export function adminListCustomers(params: components.ListCustomersRequestParams) {
	return webapi.get<components.CustomerListResponse>(`/api/admin/customers`, params);
}

/**
 * @description
 * @param req
 */
export function adminCreateCustomer(req: components.CreateCustomerRequest) {
	return webapi.post<components.CustomerInfo>(`/api/admin/customers`, req);
}

/**
 * @description
 * @param params
 */
export function adminGetCustomer(params: components.GetCustomerRequestParams, id: string) {
	return webapi.get<components.CustomerInfo>(`/api/admin/customers/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminUpdateCustomer(
	params: components.UpdateCustomerRequestParams,
	req: components.UpdateCustomerRequest,
	id: string
) {
	return webapi.put<components.CustomerInfo>(`/api/admin/customers/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function adminDeleteCustomer(params: components.DeleteCustomerRequestParams, id: string) {
	return webapi.delete<components.AnalyticsResponse>(`/api/admin/customers/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function adminListCustomerActivities(
	params: components.ListCustomerActivitiesRequestParams,
	id: string
) {
	return webapi.get<components.ListCustomerActivitiesResponse>(
		`/api/admin/customers/${id}/activities`,
		params
	);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminCancelSubscription(
	params: components.AdminCancelSubscriptionRequestParams,
	req: components.AdminCancelSubscriptionRequest,
	id: string
) {
	return webapi.post<components.AdminCancelSubscriptionResponse>(
		`/api/admin/customers/${id}/cancel-subscription`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminApplyCredit(
	params: components.ApplyCreditRequestParams,
	req: components.ApplyCreditRequest,
	id: string
) {
	return webapi.post<components.ApplyCreditResponse>(
		`/api/admin/customers/${id}/credit`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 */
export function adminGetCustomerDetail(params: components.GetCustomerRequestParams, id: string) {
	return webapi.get<components.CustomerDetailResponse>(`/api/admin/customers/${id}/detail`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminSendCustomerEmail(
	params: components.SendCustomerEmailRequestParams,
	req: components.SendCustomerEmailRequest,
	id: string
) {
	return webapi.post<components.SendCustomerEmailResponse>(
		`/api/admin/customers/${id}/email`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 */
export function adminListCustomerEmails(
	params: components.ListCustomerEmailsRequestParams,
	id: string
) {
	return webapi.get<components.ListCustomerEmailsResponse>(
		`/api/admin/customers/${id}/emails`,
		params
	);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminExtendSubscription(
	params: components.AdminExtendSubscriptionRequestParams,
	req: components.AdminExtendSubscriptionRequest,
	id: string
) {
	return webapi.post<components.AdminExtendSubscriptionResponse>(
		`/api/admin/customers/${id}/extend-subscription`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminAddCustomerNote(
	params: components.AddCustomerNoteRequestParams,
	req: components.AddCustomerNoteRequest,
	id: string
) {
	return webapi.post<components.AddCustomerNoteResponse>(
		`/api/admin/customers/${id}/notes`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 */
export function adminListCustomerNotes(
	params: components.ListCustomerNotesRequestParams,
	id: string
) {
	return webapi.get<components.ListCustomerNotesResponse>(
		`/api/admin/customers/${id}/notes`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function adminSendPasswordReset(
	params: components.SendPasswordResetRequestParams,
	id: string
) {
	return webapi.post<components.SendPasswordResetResponse>(
		`/api/admin/customers/${id}/password-reset`,
		params
	);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminRefundPayment(
	params: components.RefundPaymentRequestParams,
	req: components.RefundPaymentRequest,
	id: string
) {
	return webapi.post<components.RefundPaymentResponse>(
		`/api/admin/customers/${id}/refund`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 */
export function adminListCustomerTickets(
	params: components.ListCustomerTicketsRequestParams,
	id: string
) {
	return webapi.get<components.ListCustomerTicketsResponse>(
		`/api/admin/customers/${id}/tickets`,
		params
	);
}

/**
 * @description
 */
export function adminGetCustomerStats() {
	return webapi.get<components.CustomerStatsResponse>(`/api/admin/customers/stats`);
}

/**
 * @description
 * @param params
 */
export function listEmailDesigns(params: components.ListEmailDesignsRequestParams) {
	return webapi.get<components.ListEmailDesignsResponse>(`/api/admin/email-designs`, params);
}

/**
 * @description
 * @param req
 */
export function createEmailDesign(req: components.CreateEmailDesignRequest) {
	return webapi.post<components.EmailDesignInfo>(`/api/admin/email-designs`, req);
}

/**
 * @description
 * @param params
 */
export function getEmailDesign(params: components.GetEmailDesignRequestParams, id: string) {
	return webapi.get<components.EmailDesignInfo>(`/api/admin/email-designs/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateEmailDesign(
	params: components.UpdateEmailDesignRequestParams,
	req: components.UpdateEmailDesignRequest,
	id: string
) {
	return webapi.put<components.EmailDesignInfo>(`/api/admin/email-designs/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteEmailDesign(params: components.DeleteEmailDesignRequestParams, id: string) {
	return webapi.delete<components.AnalyticsResponse>(`/api/admin/email-designs/${id}`, params);
}

/**
 * @description
 */
export function listEmailLists() {
	return webapi.get<components.ListEmailListsResponse>(`/api/admin/email-lists`);
}

/**
 * @description
 * @param req
 */
export function createEmailList(req: components.CreateEmailListRequest) {
	return webapi.post<components.EmailListInfo>(`/api/admin/email-lists`, req);
}

/**
 * @description
 * @param params
 */
export function getEmailList(params: components.GetEmailListRequestParams, id: string) {
	return webapi.get<components.EmailListInfo>(`/api/admin/email-lists/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateEmailList(
	params: components.UpdateEmailListRequestParams,
	req: components.UpdateEmailListRequest,
	id: string
) {
	return webapi.put<components.EmailListInfo>(`/api/admin/email-lists/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteEmailList(params: components.DeleteEmailListRequestParams, id: string) {
	return webapi.delete<components.AnalyticsResponse>(`/api/admin/email-lists/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function getEmailListStats(params: components.GetEmailListRequestParams, id: string) {
	return webapi.get<components.EmailListStatsResponse>(
		`/api/admin/email-lists/${id}/stats`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function listEmailListSubscribers(
	params: components.ListSubscribersRequestParams,
	id: string
) {
	return webapi.get<components.ListSubscribersResponse>(
		`/api/admin/email-lists/${id}/subscribers`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function removeEmailListSubscriber(
	params: components.RemoveSubscriberRequestParams,
	id: string,
	subscriberId: string
) {
	return webapi.delete<components.Response>(
		`/api/admin/email-lists/${id}/subscribers/${subscriberId}`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function getEmailDashboardStats(params: components.EmailDashboardStatsRequestParams) {
	return webapi.get<components.EmailDashboardStatsResponse>(`/api/admin/email-stats`, params);
}

/**
 * @description
 * @param params
 */
export function listFunnelSteps(params: components.ListFunnelStepsRequestParams) {
	return webapi.get<components.ListFunnelStepsResponse>(`/api/admin/funnel-steps`, params);
}

/**
 * @description
 * @param req
 */
export function createFunnelStep(req: components.CreateFunnelStepRequest) {
	return webapi.post<components.FunnelStepInfo>(`/api/admin/funnel-steps`, req);
}

/**
 * @description
 * @param params
 */
export function getFunnelStep(params: components.GetFunnelStepRequestParams, id: string) {
	return webapi.get<components.FunnelStepInfo>(`/api/admin/funnel-steps/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateFunnelStep(
	params: components.UpdateFunnelStepRequestParams,
	req: components.UpdateFunnelStepRequest,
	id: string
) {
	return webapi.put<components.FunnelStepInfo>(`/api/admin/funnel-steps/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteFunnelStep(params: components.DeleteFunnelStepRequestParams, id: string) {
	return webapi.delete<components.AnalyticsResponse>(`/api/admin/funnel-steps/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function getFunnelAttribution(params: components.FunnelAttributionRequestParams) {
	return webapi.get<components.FunnelAttributionResponse>(`/api/admin/funnel/attribution`, params);
}

/**
 * @description
 * @param params
 */
export function getFunnelStats(params: components.FunnelStatsRequestParams) {
	return webapi.get<components.FunnelStatsResponse>(`/api/admin/funnel/stats`, params);
}

/**
 * @description
 * @param params
 */
export function getAllFunnelStats(params: components.AllFunnelStatsRequestParams) {
	return webapi.get<components.AllFunnelStatsResponse>(`/api/admin/funnel/stats/all`, params);
}

/**
 * @description
 */
export function listFunnels() {
	return webapi.get<components.ListFunnelsResponse>(`/api/admin/funnels`);
}

/**
 * @description
 * @param req
 */
export function createFunnel(req: components.CreateFunnelRequest) {
	return webapi.post<components.FunnelInfo>(`/api/admin/funnels`, req);
}

/**
 * @description
 * @param params
 */
export function listFunnelSubscribers(
	params: components.ListFunnelSubscribersRequestParams,
	funnel_id: number
) {
	return webapi.get<components.FunnelSubscribersResponse>(
		`/api/admin/funnels/${funnel_id}/subscribers`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function getFunnel(params: components.GetFunnelRequestParams, id: string) {
	return webapi.get<components.FunnelInfo>(`/api/admin/funnels/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateFunnel(
	params: components.UpdateFunnelRequestParams,
	req: components.UpdateFunnelRequest,
	id: string
) {
	return webapi.put<components.FunnelInfo>(`/api/admin/funnels/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteFunnel(params: components.DeleteFunnelRequestParams, id: string) {
	return webapi.delete<components.AnalyticsResponse>(`/api/admin/funnels/${id}`, params);
}

/**
 * @description
 */
export function listLeads() {
	return webapi.get<components.LeadListResponse>(`/api/admin/leads`);
}

/**
 * @description
 * @param params
 */
export function getLead(params: components.GetLeadRequestParams, id: string) {
	return webapi.get<components.LeadInfo>(`/api/admin/leads/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateLead(
	params: components.UpdateLeadRequestParams,
	req: components.UpdateLeadRequest,
	id: string
) {
	return webapi.put<components.LeadInfo>(`/api/admin/leads/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteLead(params: components.DeleteLeadRequestParams, id: string) {
	return webapi.delete<components.AnalyticsResponse>(`/api/admin/leads/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function assignLead(
	params: components.AssignLeadRequestParams,
	req: components.AssignLeadRequest,
	id: string
) {
	return webapi.put<components.LeadInfo>(`/api/admin/leads/${id}/assign`, params, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateLeadStatus(
	params: components.UpdateLeadStatusRequestParams,
	req: components.UpdateLeadStatusRequest,
	id: string
) {
	return webapi.put<components.LeadInfo>(`/api/admin/leads/${id}/status`, params, req);
}

/**
 * @description
 * @param params
 */
export function getLeadTags(params: components.LeadTagsRequestParams, id: string) {
	return webapi.get<components.LeadTagsResponse>(`/api/admin/leads/${id}/tags`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function addLeadTag(
	params: components.AddLeadTagRequestParams,
	req: components.AddLeadTagRequest,
	id: string
) {
	return webapi.post<components.AnalyticsResponse>(`/api/admin/leads/${id}/tags`, params, req);
}

/**
 * @description
 * @param params
 */
export function removeLeadTag(
	params: components.RemoveLeadTagRequestParams,
	id: string,
	tag: string
) {
	return webapi.delete<components.AnalyticsResponse>(`/api/admin/leads/${id}/tags/${tag}`, params);
}

/**
 * @description
 * @param params
 */
export function listLLMModels(params: components.ListLLMModelsRequestParams) {
	return webapi.get<components.ListLLMModelsResponse>(`/api/admin/llm/models`, params);
}

/**
 * @description
 * @param req
 */
export function createLLMModel(req: components.CreateLLMModelRequest) {
	return webapi.post<components.LLMModelPricingInfo>(`/api/admin/llm/models`, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateLLMModel(
	params: components.UpdateLLMModelRequestParams,
	req: components.UpdateLLMModelRequest,
	id: string
) {
	return webapi.put<components.LLMModelPricingInfo>(`/api/admin/llm/models/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteLLMModel(params: components.DeleteLLMModelRequestParams, id: string) {
	return webapi.delete<null>(`/api/admin/llm/models/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function listLLMRequests(params: components.ListLLMRequestsRequestParams) {
	return webapi.get<components.ListLLMRequestsResponse>(`/api/admin/llm/requests`, params);
}

/**
 * @description
 * @param params
 */
export function listOrders(params: components.ListOrdersRequestParams) {
	return webapi.get<components.ListOrdersResponse>(`/api/admin/orders`, params);
}

/**
 * @description
 * @param params
 */
export function getOrder(params: components.GetOrderRequestParams, id: string) {
	return webapi.get<components.OrderInfo>(`/api/admin/orders/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function refundOrder(
	params: components.RefundOrderRequestParams,
	req: components.RefundOrderRequest,
	id: string
) {
	return webapi.post<components.OrderInfo>(`/api/admin/orders/${id}/refund`, params, req);
}

/**
 * @description
 */
export function listOrganizations() {
	return webapi.get<components.OrgListResponse>(`/api/admin/organizations`);
}

/**
 * @description
 * @param req
 */
export function createOrganization(req: components.CreateOrgRequest) {
	return webapi.post<components.OrgInfo>(`/api/admin/organizations`, req);
}

/**
 * @description
 * @param params
 */
export function getOrganization(params: components.GetOrgRequestParams, id: string) {
	return webapi.get<components.OrgInfo>(`/api/admin/organizations/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateOrganization(
	params: components.UpdateOrgRequestParams,
	req: components.UpdateOrgRequest,
	id: string
) {
	return webapi.put<components.OrgInfo>(`/api/admin/organizations/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteOrganization(params: components.DeleteOrgRequestParams, id: string) {
	return webapi.delete<null>(`/api/admin/organizations/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function getDashboardStats(params: components.GetDashboardStatsRequestParams, id: string) {
	return webapi.get<components.DashboardStatsResponse>(
		`/api/admin/organizations/${id}/dashboard-stats`,
		params
	);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateOrgEmailSettings(
	params: components.UpdateOrgEmailSettingsRequestParams,
	req: components.UpdateOrgEmailSettingsRequest,
	id: string
) {
	return webapi.put<components.OrgInfo>(`/api/admin/organizations/${id}/email`, params, req);
}

/**
 * @description
 * @param params
 */
export function getOrgLLMCredentials(
	params: components.GetOrgLLMCredentialsRequestParams,
	id: string
) {
	return webapi.get<components.OrgLLMCredentialsResponse>(
		`/api/admin/organizations/${id}/llm-credentials`,
		params
	);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateOrgLLMCredentials(
	params: components.UpdateOrgLLMCredentialsRequestParams,
	req: components.UpdateOrgLLMCredentialsRequest,
	id: string
) {
	return webapi.put<components.OrgLLMCredentialsResponse>(
		`/api/admin/organizations/${id}/llm-credentials`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateOrgLLMPricing(
	params: components.UpdateOrgLLMPricingRequestParams,
	req: components.UpdateOrgLLMPricingRequest,
	id: string
) {
	return webapi.put<components.OrgInfo>(`/api/admin/organizations/${id}/llm-pricing`, params, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function testOrgLLMConnection(
	params: components.TestOrgLLMConnectionRequestParams,
	req: components.TestOrgLLMConnectionRequest,
	id: string
) {
	return webapi.post<components.TestOrgLLMConnectionResponse>(
		`/api/admin/organizations/${id}/llm-test`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 */
export function regenerateOrgApiKey(params: components.RegenerateApiKeyRequestParams, id: string) {
	return webapi.post<components.OrgInfo>(`/api/admin/organizations/${id}/regenerate-key`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateOrgStripe(
	params: components.UpdateOrgStripeRequestParams,
	req: components.UpdateOrgStripeRequest,
	id: string
) {
	return webapi.put<components.OrgInfo>(`/api/admin/organizations/${id}/stripe`, params, req);
}

/**
 * @description
 * @param params
 */
export function getOrganizationBySlug(params: components.GetOrgBySlugRequestParams, slug: string) {
	return webapi.get<components.OrgInfo>(`/api/admin/organizations/slug/${slug}`, params);
}

/**
 * @description
 * @param params
 */
export function adminListAutomationLog(params: components.ListAutomationLogRequestParams) {
	return webapi.get<components.ListAutomationLogResponse>(`/api/admin/automation-log`, params);
}

/**
 * @description
 * @param params
 */
export function adminListRuleTemplates(params: components.ListRuleTemplatesRequestParams) {
	return webapi.get<components.ListRuleTemplatesResponse>(`/api/admin/rule-templates`, params);
}

/**
 * @description
 * @param params
 */
export function adminListRules(params: components.ListRulesRequestParams) {
	return webapi.get<components.ListRulesResponse>(`/api/admin/rules`, params);
}

/**
 * @description
 * @param req
 */
export function adminCreateRule(req: components.CreateRuleRequest) {
	return webapi.post<components.RuleInfo>(`/api/admin/rules`, req);
}

/**
 * @description
 * @param params
 */
export function adminGetRule(params: components.GetRuleRequestParams, id: string) {
	return webapi.get<components.RuleInfo>(`/api/admin/rules/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminUpdateRule(
	params: components.UpdateRuleRequestParams,
	req: components.UpdateRuleRequest,
	id: string
) {
	return webapi.put<components.RuleInfo>(`/api/admin/rules/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function adminDeleteRule(params: components.DeleteRuleRequestParams, id: string) {
	return webapi.delete<components.EmptyResponse>(`/api/admin/rules/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminTestRule(
	params: components.TestRuleRequestParams,
	req: components.TestRuleRequest,
	id: string
) {
	return webapi.post<components.TestRuleResponse>(`/api/admin/rules/${id}/test`, params, req);
}

/**
 * @description
 * @param params
 */
export function adminToggleRule(params: components.ToggleRuleRequestParams, id: string) {
	return webapi.post<components.RuleInfo>(`/api/admin/rules/${id}/toggle`, params);
}

/**
 * @description
 */
export function adminRestoreDefaultRules() {
	return webapi.post<components.RestoreDefaultRulesResponse>(`/api/admin/rules/restore-defaults`);
}

/**
 * @description
 */
export function adminGetRulesSchema() {
	return webapi.get<components.RulesSchemaResponse>(`/api/admin/rules/schema`);
}

/**
 * @description
 * @param params
 */
export function listEmailQueue(params: components.EmailQueueListRequestParams) {
	return webapi.get<components.EmailQueueListResponse>(`/api/admin/email-queue`, params);
}

/**
 * @description
 * @param params
 */
export function cancelEmail(params: components.CancelEmailRequestParams, id: string) {
	return webapi.post<components.AnalyticsResponse>(`/api/admin/email-queue/${id}/cancel`, params);
}

/**
 * @description
 * @param req
 */
export function createEntryRule(req: components.CreateEntryRuleRequest) {
	return webapi.post<components.EntryRuleResponse>(`/api/admin/entry-rules`, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateEntryRule(
	params: components.UpdateEntryRuleRequestParams,
	req: components.UpdateEntryRuleRequest,
	id: string
) {
	return webapi.put<components.EntryRuleResponse>(`/api/admin/entry-rules/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteEntryRule(params: components.DeleteEntryRuleRequestParams, id: string) {
	return webapi.delete<components.AnalyticsResponse>(`/api/admin/entry-rules/${id}`, params);
}

/**
 * @description
 */
export function listSequences() {
	return webapi.get<components.SequenceListResponse>(`/api/admin/sequences`);
}

/**
 * @description
 * @param req
 */
export function createSequence(req: components.CreateSequenceRequest) {
	return webapi.post<components.SequenceInfo>(`/api/admin/sequences`, req);
}

/**
 * @description
 * @param params
 */
export function getSequence(params: components.GetSequenceRequestParams, id: string) {
	return webapi.get<components.SequenceDetailResponse>(`/api/admin/sequences/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateSequence(
	params: components.UpdateSequenceRequestParams,
	req: components.UpdateSequenceRequest,
	id: string
) {
	return webapi.put<components.SequenceInfo>(`/api/admin/sequences/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteSequence(params: components.GetSequenceRequestParams, id: string) {
	return webapi.delete<components.AnalyticsResponse>(`/api/admin/sequences/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function getSequenceStats(params: components.SequenceStatsRequestParams, id: string) {
	return webapi.get<components.SequenceStatsResponse>(`/api/admin/sequences/${id}/stats`, params);
}

/**
 * @description
 * @param req
 */
export function createTemplate(req: components.CreateTemplateRequest) {
	return webapi.post<components.TemplateInfo>(`/api/admin/templates`, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateTemplate(
	params: components.UpdateTemplateRequestParams,
	req: components.UpdateTemplateRequest,
	id: string
) {
	return webapi.put<components.TemplateInfo>(`/api/admin/templates/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteTemplate(params: components.DeleteTemplateRequestParams, id: string) {
	return webapi.delete<components.AnalyticsResponse>(`/api/admin/templates/${id}`, params);
}

/**
 * @description
 */
export function getPlatformSettings() {
	return webapi.get<components.GetPlatformSettingsResponse>(`/api/admin/settings`);
}

/**
 * @description
 * @param params
 */
export function getPlatformSettingsByCategory(
	params: components.GetPlatformSettingsByCategoryRequestParams,
	category: string
) {
	return webapi.get<components.GetPlatformSettingsResponse>(
		`/api/admin/settings/${category}`,
		params
	);
}

/**
 * @description
 * @param req
 */
export function updateAISettings(req: components.UpdateAISettingsRequest) {
	return webapi.put<components.UpdateSettingsResponse>(`/api/admin/settings/ai`, req);
}

/**
 * @description
 * @param req
 */
export function testAIConnection(req: components.TestAIConnectionRequest) {
	return webapi.post<components.TestAIConnectionResponse>(`/api/admin/settings/ai/test`, req);
}

/**
 * @description
 * @param req
 */
export function updateEmailSettings(req: components.UpdateEmailSettingsRequest) {
	return webapi.put<components.UpdateSettingsResponse>(`/api/admin/settings/email`, req);
}

/**
 * @description
 * @param req
 */
export function updateGoogleSettings(req: components.UpdateGoogleSettingsRequest) {
	return webapi.put<components.UpdateSettingsResponse>(`/api/admin/settings/google`, req);
}

/**
 * @description
 * @param params
 */
export function blockSubscriber(params: components.SubscriberActionRequestParams, id: string) {
	return webapi.post<components.FunnelSubscriberInfo>(`/api/admin/subscribers/${id}/block`, params);
}

/**
 * @description
 * @param params
 */
export function resubscribeSubscriber(
	params: components.SubscriberActionRequestParams,
	id: string
) {
	return webapi.post<components.FunnelSubscriberInfo>(
		`/api/admin/subscribers/${id}/resubscribe`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function unblockSubscriber(params: components.SubscriberActionRequestParams, id: string) {
	return webapi.post<components.FunnelSubscriberInfo>(
		`/api/admin/subscribers/${id}/unblock`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function unsubscribeSubscriber(
	params: components.SubscriberActionRequestParams,
	id: string
) {
	return webapi.post<components.FunnelSubscriberInfo>(
		`/api/admin/subscribers/${id}/unsubscribe`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function verifySubscriber(params: components.SubscriberActionRequestParams, id: string) {
	return webapi.post<components.FunnelSubscriberInfo>(
		`/api/admin/subscribers/${id}/verify`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function adminListTickets(params: components.ListTicketsRequestParams) {
	return webapi.get<components.ListTicketsResponse>(`/api/admin/tickets`, params);
}

/**
 * @description
 * @param req
 */
export function adminCreateTicket(req: components.CreateTicketRequest) {
	return webapi.post<components.SupportTicketInfo>(`/api/admin/tickets`, req);
}

/**
 * @description
 * @param params
 */
export function adminGetTicket(params: components.GetTicketRequestParams, id: string) {
	return webapi.get<components.TicketDetailResponse>(`/api/admin/tickets/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminUpdateTicket(
	params: components.UpdateTicketRequestParams,
	req: components.UpdateTicketRequest,
	id: string
) {
	return webapi.put<components.SupportTicketInfo>(`/api/admin/tickets/${id}`, params, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminAddTicketMessage(
	params: components.AddTicketMessageRequestParams,
	req: components.AddTicketMessageRequest,
	id: string
) {
	return webapi.post<components.AddTicketMessageResponse>(
		`/api/admin/tickets/${id}/messages`,
		params,
		req
	);
}

/**
 * @description
 */
export function adminGetTicketStats() {
	return webapi.get<components.TicketStatsResponse>(`/api/admin/tickets/stats`);
}

/**
 * @description
 */
export function listTransactionalEmails() {
	return webapi.get<components.ListTransactionalEmailsResponse>(`/api/admin/transactional-emails`);
}

/**
 * @description
 * @param req
 */
export function createTransactionalEmail(req: components.CreateTransactionalEmailRequest) {
	return webapi.post<components.TransactionalEmailInfo>(`/api/admin/transactional-emails`, req);
}

/**
 * @description
 * @param params
 */
export function getTransactionalEmail(
	params: components.GetTransactionalEmailRequestParams,
	id: string
) {
	return webapi.get<components.TransactionalEmailInfo>(
		`/api/admin/transactional-emails/${id}`,
		params
	);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateTransactionalEmail(
	params: components.UpdateTransactionalEmailRequestParams,
	req: components.UpdateTransactionalEmailRequest,
	id: string
) {
	return webapi.put<components.TransactionalEmailInfo>(
		`/api/admin/transactional-emails/${id}`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 */
export function deleteTransactionalEmail(
	params: components.DeleteTransactionalEmailRequestParams,
	id: string
) {
	return webapi.delete<components.AnalyticsResponse>(
		`/api/admin/transactional-emails/${id}`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function getTransactionalEmailStats(
	params: components.GetTransactionalEmailRequestParams,
	id: string
) {
	return webapi.get<components.TransactionalStatsResponse>(
		`/api/admin/transactional-emails/${id}/stats`,
		params
	);
}

/**
 * @description
 */
export function listUsers() {
	return webapi.get<components.UserListResponse>(`/api/admin/users`);
}

/**
 * @description
 * @param req
 */
export function createUser(req: components.CreateUserRequest) {
	return webapi.post<components.CreateUserResponse>(`/api/admin/users`, req);
}

/**
 * @description
 * @param params
 */
export function getUser(params: components.GetUserRequestParams, id: string) {
	return webapi.get<components.UserInfo>(`/api/admin/users/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateUser(
	params: components.UpdateUserRequestParams,
	req: components.UpdateUserRequest,
	id: string
) {
	return webapi.put<components.UserInfo>(`/api/admin/users/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteUser(params: components.DeleteUserRequestParams, id: string) {
	return webapi.delete<components.AnalyticsResponse>(`/api/admin/users/${id}`, params);
}

/**
 * @description
 */
export function adminListWebhooks() {
	return webapi.get<components.ListWebhooksResponse>(`/api/admin/webhooks`);
}

/**
 * @description
 * @param req
 */
export function adminCreateWebhook(req: components.RegisterWebhookRequest) {
	return webapi.post<components.RegisterWebhookResponse>(`/api/admin/webhooks`, req);
}

/**
 * @description
 * @param params
 */
export function adminGetWebhook(params: components.GetWebhookRequestParams, id: string) {
	return webapi.get<components.WebhookInfo>(`/api/admin/webhooks/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function adminUpdateWebhook(
	params: components.UpdateWebhookRequestParams,
	req: components.UpdateWebhookRequest,
	id: string
) {
	return webapi.put<components.WebhookInfo>(`/api/admin/webhooks/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function adminDeleteWebhook(params: components.DeleteWebhookRequestParams, id: string) {
	return webapi.delete<components.Response>(`/api/admin/webhooks/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function adminListWebhookLogs(params: components.ListWebhookLogsRequestParams, id: string) {
	return webapi.get<components.ListWebhookLogsResponse>(`/api/admin/webhooks/${id}/logs`, params);
}

/**
 * @description
 * @param params
 */
export function adminTestWebhook(params: components.TestWebhookRequestParams, id: string) {
	return webapi.post<components.TestWebhookResponse>(`/api/admin/webhooks/${id}/test`, params);
}

/**
 * @description
 * @param params
 */
export function listWorkshopRegistrations(
	params: components.ListWorkshopRegistrationsRequestParams
) {
	return webapi.get<components.ListWorkshopRegistrationsResponse>(
		`/api/admin/workshop-registrations`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function getWorkshopRegistration(
	params: components.GetWorkshopRegistrationRequestParams,
	id: string
) {
	return webapi.get<components.WorkshopRegistrationInfo>(
		`/api/admin/workshop-registrations/${id}`,
		params
	);
}

/**
 * @description
 * @param params
 * @param req
 */
export function cancelWorkshopRegistration(
	params: components.CancelWorkshopRegistrationRequestParams,
	req: components.CancelWorkshopRegistrationRequest,
	id: string
) {
	return webapi.post<components.WorkshopRegistrationInfo>(
		`/api/admin/workshop-registrations/${id}/cancel`,
		params,
		req
	);
}

/**
 * @description
 * @param params
 */
export function listWorkshopEvents(params: components.ListWorkshopEventsRequestParams) {
	return webapi.get<components.ListWorkshopEventsResponse>(`/api/admin/workshops`, params);
}

/**
 * @description
 * @param req
 */
export function createWorkshopEvent(req: components.CreateWorkshopEventRequest) {
	return webapi.post<components.WorkshopEventInfo>(`/api/admin/workshops`, req);
}

/**
 * @description
 * @param params
 */
export function getWorkshopEvent(params: components.GetWorkshopEventRequestParams, id: string) {
	return webapi.get<components.WorkshopEventInfo>(`/api/admin/workshops/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateWorkshopEvent(
	params: components.UpdateWorkshopEventRequestParams,
	req: components.UpdateWorkshopEventRequest,
	id: string
) {
	return webapi.put<components.WorkshopEventInfo>(`/api/admin/workshops/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteWorkshopEvent(
	params: components.DeleteWorkshopEventRequestParams,
	id: string
) {
	return webapi.delete<components.AnalyticsResponse>(`/api/admin/workshops/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateWorkshopStatus(
	params: components.UpdateWorkshopStatusRequestParams,
	req: components.UpdateWorkshopStatusRequest,
	id: string
) {
	return webapi.put<components.WorkshopEventInfo>(`/api/admin/workshops/${id}/status`, params, req);
}

/**
 * @description "Request password reset"
 * @param req
 */
export function forgotPassword(req: components.ForgotPasswordRequest) {
	return webapi.post<components.ForgotPasswordResponse>(`/api/v1/auth/forgot-password`, req);
}

/**
 * @description "User login"
 * @param req
 */
export function login(req: components.LoginRequest) {
	return webapi.post<components.LoginResponse>(`/api/v1/auth/login`, req);
}

/**
 * @description "Refresh access token using refresh token"
 * @param req
 */
export function refreshToken(req: components.RefreshTokenRequest) {
	return webapi.post<components.LoginResponse>(`/api/v1/auth/refresh`, req);
}

/**
 * @description "User registration"
 * @param req
 */
export function register(req: components.RegisterRequest) {
	return webapi.post<components.RegisterResponse>(`/api/v1/auth/register`, req);
}

/**
 * @description "Reset password with token"
 * @param req
 */
export function resetPassword(req: components.ResetPasswordRequest) {
	return webapi.post<components.ResetPasswordResponse>(`/api/v1/auth/reset-password`, req);
}

/**
 * @description "Validate invitation token and get user details"
 * @param params
 */
export function validateInvitation(params: components.ValidateInvitationRequestParams) {
	return webapi.get<components.ValidateInvitationResponse>(
		`/api/v1/auth/validate-invitation`,
		params
	);
}

/**
 * @description "Verify email address"
 * @param params
 */
export function verifyEmail(params: components.VerifyEmailRequestParams) {
	return webapi.get<components.VerifyEmailResponse>(`/api/v1/auth/verify-email`, params);
}

/**
 * @description
 * @param params
 */
export function getAvailability(params: components.AvailabilityRequestParams) {
	return webapi.get<components.AvailabilityResponse>(`/api/calendar/availability`, params);
}

/**
 * @description
 * @param req
 */
export function bookMeeting(req: components.BookMeetingRequest) {
	return webapi.post<components.BookMeetingResponse>(`/api/calendar/book`, req);
}

/**
 * @description
 */
export function oAuthConnect() {
	return webapi.get<components.OAuthConnectResponse>(`/api/oauth/google/connect`);
}

/**
 * @description
 */
export function oAuthDisconnect() {
	return webapi.post<components.AnalyticsResponse>(`/api/oauth/google/disconnect`);
}

/**
 * @description
 */
export function oAuthStatus() {
	return webapi.get<components.OAuthStatusResponse>(`/api/oauth/google/status`);
}

/**
 * @description
 */
export function health() {
	return webapi.get<components.HealthResponse>(`/api/v1/health`);
}

/**
 * @description
 * @param req
 */
export function createContact(req: components.ContactRequest) {
	return webapi.post<components.ContactResponse>(`/sdk/v1/contacts`, req);
}

/**
 * @description
 * @param req
 */
export function trackEvent(req: components.EventRequest) {
	return webapi.post<components.Response>(`/sdk/v1/events`, req);
}

/**
 * @description
 * @param params
 */
export function sDKGetFunnelStep(
	params: components.GetFunnelStepBySlugRequestParams,
	slug: string
) {
	return webapi.get<components.FunnelStepInfo>(`/sdk/v1/funnels/${slug}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function subscribeToList(
	params: components.SubscribeRequestParams,
	req: components.SubscribeRequest,
	slug: string
) {
	return webapi.post<components.Response>(`/sdk/v1/lists/${slug}/subscribe`, params, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function unsubscribeFromList(
	params: components.SubscribeRequestParams,
	req: components.SubscribeRequest,
	slug: string
) {
	return webapi.post<components.Response>(`/sdk/v1/lists/${slug}/unsubscribe`, params, req);
}

/**
 * @description
 * @param req
 */
export function processOffer(req: components.OfferRequest) {
	return webapi.post<components.OfferResponse>(`/sdk/v1/offers`, req);
}

/**
 * @description
 * @param req
 */
export function createOrder(req: components.OrderRequest) {
	return webapi.post<components.OrderResponse>(`/sdk/v1/orders`, req);
}

/**
 * @description
 * @param params
 */
export function sDKGetProduct(params: components.GetProductBySlugRequestParams, slug: string) {
	return webapi.get<components.BillingProductInfo>(`/sdk/v1/products/${slug}`, params);
}

/**
 * @description
 * @param params
 */
export function getQuiz(params: components.GetQuizBySlugRequestParams, slug: string) {
	return webapi.get<components.QuizInfo>(`/sdk/v1/quizzes/${slug}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function submitQuiz(
	params: components.QuizSubmitRequestParams,
	req: components.QuizSubmitRequest,
	slug: string
) {
	return webapi.post<components.QuizSubmitResponse>(`/sdk/v1/quizzes/${slug}/submit`, params, req);
}

/**
 * @description
 * @param params
 */
export function sDKGetWorkshop(params: components.GetWorkshopBySlugRequestParams, slug: string) {
	return webapi.get<components.WorkshopEventInfo>(`/sdk/v1/workshops/${slug}`, params);
}

/**
 * @description
 * @param params
 */
export function sDKGetWorkshopByProduct(
	params: components.GetWorkshopByProductSlugRequestParams,
	productSlug: string
) {
	return webapi.get<components.WorkshopEventInfo>(
		`/sdk/v1/workshops/product/${productSlug}`,
		params
	);
}

/**
 * @description
 * @param req
 */
export function sDKChangePassword(req: components.SDKChangePasswordRequest) {
	return webapi.post<components.Response>(`/sdk/v1/auth/change-password`, req);
}

/**
 * @description
 * @param req
 */
export function sDKForgotPassword(req: components.SDKForgotPasswordRequest) {
	return webapi.post<components.Response>(`/sdk/v1/auth/forgot-password`, req);
}

/**
 * @description
 * @param req
 */
export function sDKLogin(req: components.SDKLoginRequest) {
	return webapi.post<components.SDKAuthResponse>(`/sdk/v1/auth/login`, req);
}

/**
 * @description
 * @param req
 */
export function sDKRefreshToken(req: components.SDKRefreshTokenRequest) {
	return webapi.post<components.SDKAuthResponse>(`/sdk/v1/auth/refresh`, req);
}

/**
 * @description
 * @param req
 */
export function sDKRegister(req: components.SDKRegisterRequest) {
	return webapi.post<components.SDKAuthResponse>(`/sdk/v1/auth/register`, req);
}

/**
 * @description
 * @param req
 */
export function sDKResetPassword(req: components.SDKResetPasswordRequest) {
	return webapi.post<components.Response>(`/sdk/v1/auth/reset-password`, req);
}

/**
 * @description
 * @param req
 */
export function sDKVerifyEmail(req: components.SDKVerifyEmailRequest) {
	return webapi.post<components.Response>(`/sdk/v1/auth/verify-email`, req);
}

/**
 * @description
 * @param req
 */
export function createCheckoutSession(req: components.CheckoutRequest) {
	return webapi.post<components.CheckoutResponse>(`/sdk/v1/billing/checkout`, req);
}

/**
 * @description
 * @param req
 */
export function createCustomer(req: components.CustomerRequest) {
	return webapi.post<components.CustomerResponse>(`/sdk/v1/billing/customers`, req);
}

/**
 * @description
 * @param req
 */
export function getCustomerPortal(req: components.PortalRequest) {
	return webapi.post<components.PortalResponse>(`/sdk/v1/billing/portal`, req);
}

/**
 * @description
 * @param req
 */
export function createSubscription(req: components.SubscriptionRequest) {
	return webapi.post<components.SubscriptionResponse>(`/sdk/v1/billing/subscriptions`, req);
}

/**
 * @description
 * @param params
 */
export function cancelSubscription(params: components.CancelSubscriptionRequestParams, id: string) {
	return webapi.post<components.Response>(`/sdk/v1/billing/subscriptions/${id}/cancel`, params);
}

/**
 * @description
 * @param req
 */
export function recordUsage(req: components.UsageRequest) {
	return webapi.post<components.Response>(`/sdk/v1/billing/usage`, req);
}

/**
 * @description
 * @param params
 */
export function getContact(params: components.GetContactRequestParams, id: string) {
	return webapi.get<components.SDKContactInfo>(`/sdk/v1/contacts/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateContact(
	params: components.UpdateContactRequestParams,
	req: components.UpdateContactRequest,
	id: string
) {
	return webapi.put<components.SDKContactInfo>(`/sdk/v1/contacts/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function listContactActivity(
	params: components.ListContactActivityRequestParams,
	id: string
) {
	return webapi.get<components.ListContactActivityResponse>(
		`/sdk/v1/contacts/${id}/activity`,
		params
	);
}

/**
 * @description
 * @param params
 * @param req
 */
export function addContactTags(
	params: components.AddContactTagsRequestParams,
	req: components.AddContactTagsRequest,
	id: string
) {
	return webapi.post<components.Response>(`/sdk/v1/contacts/${id}/tags`, params, req);
}

/**
 * @description
 * @param params
 * @param req
 */
export function removeContactTags(
	params: components.RemoveContactTagsRequestParams,
	req: components.RemoveContactTagsRequest,
	id: string
) {
	return webapi.delete<components.Response>(`/sdk/v1/contacts/${id}/tags`, params, req);
}

/**
 * @description
 * @param req
 */
export function globalUnsubscribe(req: components.GlobalUnsubscribeRequest) {
	return webapi.post<components.Response>(`/sdk/v1/contacts/unsubscribe`, req);
}

/**
 * @description
 */
export function listSDKContentCategories() {
	return webapi.get<components.ListSDKContentCategoriesResponse>(`/sdk/v1/content/categories`);
}

/**
 * @description
 * @param params
 */
export function listSDKContentPages(params: components.ListSDKContentPagesRequestParams) {
	return webapi.get<components.ListSDKContentPagesResponse>(`/sdk/v1/content/pages`, params);
}

/**
 * @description
 * @param params
 */
export function getSDKContentPage(params: components.GetSDKContentPageRequestParams, slug: string) {
	return webapi.get<components.SDKContentPageInfo>(`/sdk/v1/content/pages/${slug}`, params);
}

/**
 * @description
 * @param params
 */
export function listSDKContentPosts(params: components.ListSDKContentPostsRequestParams) {
	return webapi.get<components.ListSDKContentPostsResponse>(`/sdk/v1/content/posts`, params);
}

/**
 * @description
 * @param params
 */
export function getSDKContentPost(params: components.GetSDKContentPostRequestParams, slug: string) {
	return webapi.get<components.SDKContentPostInfo>(`/sdk/v1/content/posts/${slug}`, params);
}

/**
 * @description
 * @param params
 */
export function getCustomerByEmail(
	params: components.GetCustomerByEmailRequestParams,
	email: string
) {
	return webapi.get<components.SDKCustomerInfo>(`/sdk/v1/customers/${email}`, params);
}

/**
 * @description
 * @param params
 */
export function listCustomerInvoices(
	params: components.ListCustomerInvoicesRequestParams,
	email: string
) {
	return webapi.get<components.ListCustomerInvoicesResponse>(
		`/sdk/v1/customers/${email}/invoices`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function listCustomerOrders(
	params: components.ListCustomerOrdersRequestParams,
	email: string
) {
	return webapi.get<components.ListCustomerOrdersResponse>(
		`/sdk/v1/customers/${email}/orders`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function listCustomerPayments(
	params: components.ListCustomerPaymentsRequestParams,
	email: string
) {
	return webapi.get<components.ListCustomerPaymentsResponse>(
		`/sdk/v1/customers/${email}/payments`,
		params
	);
}

/**
 * @description
 * @param params
 */
export function listCustomerSubscriptions(
	params: components.ListCustomerSubscriptionsRequestParams,
	email: string
) {
	return webapi.get<components.ListCustomerSubscriptionsResponse>(
		`/sdk/v1/customers/${email}/subscriptions`,
		params
	);
}

/**
 * @description
 * @param params
 * @param req
 */
export function sDKUpdateCustomer(
	params: components.SDKUpdateCustomerRequestParams,
	req: components.SDKUpdateCustomerRequest,
	id: string
) {
	return webapi.put<components.SDKCustomerInfo>(`/sdk/v1/customers/id/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function sDKDeleteCustomer(params: components.SDKDeleteCustomerRequestParams, id: string) {
	return webapi.delete<components.Response>(`/sdk/v1/customers/id/${id}`, params);
}

/**
 * @description
 * @param req
 */
export function sendEmail(req: components.SendEmailRequest) {
	return webapi.post<components.SendEmailResponse>(`/sdk/v1/emails`, req);
}

/**
 * @description
 * @param params
 */
export function getEmailStatus(params: components.GetEmailStatusRequestParams, messageId: string) {
	return webapi.get<components.EmailStatusResponse>(`/sdk/v1/emails/${messageId}`, params);
}

/**
 * @description
 * @param params
 */
export function listEmailEvents(
	params: components.ListEmailEventsRequestParams,
	messageId: string
) {
	return webapi.get<components.ListEmailEventsResponse>(
		`/sdk/v1/emails/${messageId}/events`,
		params
	);
}

/**
 * @description
 * @param req
 */
export function lLMChat(req: components.LLMChatRequest) {
	return webapi.post<components.LLMChatResponse>(`/sdk/v1/llm/chat`, req);
}

/**
 * @description
 */
export function lLMConfig() {
	return webapi.get<components.LLMConfigResponse>(`/sdk/v1/llm/config`);
}

/**
 * @description
 * @param req
 */
export function enrollInSequence(req: components.EnrollSequenceRequest) {
	return webapi.post<components.EnrollSequenceResponse>(`/sdk/v1/sequences/enroll`, req);
}

/**
 * @description
 * @param params
 */
export function getSequenceEnrollments(params: components.GetSequenceEnrollmentRequestParams) {
	return webapi.get<components.ListSequenceEnrollmentsResponse>(
		`/sdk/v1/sequences/enrollments`,
		params
	);
}

/**
 * @description
 * @param req
 */
export function pauseSequenceEnrollment(req: components.PauseSequenceRequest) {
	return webapi.post<components.Response>(`/sdk/v1/sequences/pause`, req);
}

/**
 * @description
 * @param req
 */
export function resumeSequenceEnrollment(req: components.ResumeSequenceRequest) {
	return webapi.post<components.Response>(`/sdk/v1/sequences/resume`, req);
}

/**
 * @description
 * @param req
 */
export function unenrollFromSequence(req: components.UnenrollSequenceRequest) {
	return webapi.post<components.Response>(`/sdk/v1/sequences/unenroll`, req);
}

/**
 * @description
 */
export function listSDKAuthors() {
	return webapi.get<components.ListSDKAuthorsResponse>(`/sdk/v1/site/authors`);
}

/**
 * @description
 * @param params
 */
export function getSDKAuthor(params: components.GetSDKAuthorRequestParams, id: string) {
	return webapi.get<components.SDKAuthorInfo>(`/sdk/v1/site/authors/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function listSDKNavigationMenus(params: components.ListSDKNavigationMenusRequestParams) {
	return webapi.get<components.ListSDKNavigationMenusResponse>(`/sdk/v1/site/menus`, params);
}

/**
 * @description
 * @param params
 */
export function getSDKNavigationMenu(
	params: components.GetSDKNavigationMenuRequestParams,
	slug: string
) {
	return webapi.get<components.SDKNavigationMenu>(`/sdk/v1/site/menus/${slug}`, params);
}

/**
 * @description
 */
export function getSDKSiteSettings() {
	return webapi.get<components.SDKSiteSettings>(`/sdk/v1/site/settings`);
}

/**
 * @description
 * @param params
 */
export function getContactStats(params: components.GetContactStatsRequestParams) {
	return webapi.get<components.GetContactStatsResponse>(`/sdk/v1/stats/contacts`, params);
}

/**
 * @description
 * @param params
 */
export function getEmailStats(params: components.GetEmailStatsRequestParams) {
	return webapi.get<components.GetEmailStatsResponse>(`/sdk/v1/stats/emails`, params);
}

/**
 * @description
 * @param params
 */
export function getStatsOverview(params: components.GetStatsOverviewRequestParams) {
	return webapi.get<components.StatsOverviewResponse>(`/sdk/v1/stats/overview`, params);
}

/**
 * @description
 * @param params
 */
export function getRevenueStats(params: components.GetRevenueStatsRequestParams) {
	return webapi.get<components.GetRevenueStatsResponse>(`/sdk/v1/stats/revenue`, params);
}

/**
 * @description
 * @param req
 */
export function trackClick(req: components.TrackClickRequest) {
	return webapi.post<components.Response>(`/sdk/v1/tracking/click`, req);
}

/**
 * @description
 * @param req
 */
export function trackConfirm(req: components.TrackConfirmRequest) {
	return webapi.post<components.TrackConfirmResponse>(`/sdk/v1/tracking/confirm`, req);
}

/**
 * @description
 * @param req
 */
export function trackOpen(req: components.TrackOpenRequest) {
	return webapi.post<components.Response>(`/sdk/v1/tracking/open`, req);
}

/**
 * @description
 * @param req
 */
export function trackUnsubscribe(req: components.TrackUnsubscribeRequest) {
	return webapi.post<components.Response>(`/sdk/v1/tracking/unsubscribe`, req);
}

/**
 * @description
 * @param req
 */
export function registerWebhook(req: components.RegisterWebhookRequest) {
	return webapi.post<components.RegisterWebhookResponse>(`/sdk/v1/webhooks`, req);
}

/**
 * @description
 */
export function listWebhooks() {
	return webapi.get<components.ListWebhooksResponse>(`/sdk/v1/webhooks`);
}

/**
 * @description
 * @param params
 */
export function getWebhook(params: components.GetWebhookRequestParams, id: string) {
	return webapi.get<components.WebhookInfo>(`/sdk/v1/webhooks/${id}`, params);
}

/**
 * @description
 * @param params
 * @param req
 */
export function updateWebhook(
	params: components.UpdateWebhookRequestParams,
	req: components.UpdateWebhookRequest,
	id: string
) {
	return webapi.put<components.WebhookInfo>(`/sdk/v1/webhooks/${id}`, params, req);
}

/**
 * @description
 * @param params
 */
export function deleteWebhook(params: components.DeleteWebhookRequestParams, id: string) {
	return webapi.delete<components.Response>(`/sdk/v1/webhooks/${id}`, params);
}

/**
 * @description
 * @param params
 */
export function listWebhookLogs(params: components.ListWebhookLogsRequestParams, id: string) {
	return webapi.get<components.ListWebhookLogsResponse>(`/sdk/v1/webhooks/${id}/logs`, params);
}

/**
 * @description
 * @param params
 */
export function testWebhook(params: components.TestWebhookRequestParams, id: string) {
	return webapi.post<components.TestWebhookResponse>(`/sdk/v1/webhooks/${id}/test`, params);
}
