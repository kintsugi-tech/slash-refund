/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

/**
* `Any` contains an arbitrary serialized protocol buffer message along with a
URL that describes the type of the serialized message.

Protobuf library provides support to pack/unpack Any values in the form
of utility functions or additional generated methods of the Any type.

Example 1: Pack and unpack a message in C++.

    Foo foo = ...;
    Any any;
    any.PackFrom(foo);
    ...
    if (any.UnpackTo(&foo)) {
      ...
    }

Example 2: Pack and unpack a message in Java.

    Foo foo = ...;
    Any any = Any.pack(foo);
    ...
    if (any.is(Foo.class)) {
      foo = any.unpack(Foo.class);
    }

 Example 3: Pack and unpack a message in Python.

    foo = Foo(...)
    any = Any()
    any.Pack(foo)
    ...
    if any.Is(Foo.DESCRIPTOR):
      any.Unpack(foo)
      ...

 Example 4: Pack and unpack a message in Go

     foo := &pb.Foo{...}
     any, err := anypb.New(foo)
     if err != nil {
       ...
     }
     ...
     foo := &pb.Foo{}
     if err := any.UnmarshalTo(foo); err != nil {
       ...
     }

The pack methods provided by protobuf library will by default use
'type.googleapis.com/full.type.name' as the type URL and the unpack
methods only use the fully qualified type name after the last '/'
in the type URL, for example "foo.bar.com/x/y.z" will yield type
name "y.z".


JSON
====
The JSON representation of an `Any` value uses the regular
representation of the deserialized, embedded message, with an
additional field `@type` which contains the type URL. Example:

    package google.profile;
    message Person {
      string first_name = 1;
      string last_name = 2;
    }

    {
      "@type": "type.googleapis.com/google.profile.Person",
      "firstName": <string>,
      "lastName": <string>
    }

If the embedded message type is well-known and has a custom JSON
representation, that representation will be embedded adding a field
`value` which holds the custom JSON in addition to the `@type`
field. Example (for message [google.protobuf.Duration][]):

    {
      "@type": "type.googleapis.com/google.protobuf.Duration",
      "value": "1.212s"
    }
*/
export interface ProtobufAny {
  /**
   * A URL/resource name that uniquely identifies the type of the serialized
   * protocol buffer message. This string must contain at least
   * one "/" character. The last segment of the URL's path must represent
   * the fully qualified name of the type (as in
   * `path/google.protobuf.Duration`). The name should be in a canonical form
   * (e.g., leading "." is not accepted).
   *
   * In practice, teams usually precompile into the binary all types that they
   * expect it to use in the context of Any. However, for URLs which use the
   * scheme `http`, `https`, or no scheme, one can optionally set up a type
   * server that maps type URLs to message definitions as follows:
   * * If no scheme is provided, `https` is assumed.
   * * An HTTP GET on the URL must yield a [google.protobuf.Type][]
   *   value in binary format, or produce an error.
   * * Applications are allowed to cache lookup results based on the
   *   URL, or have them precompiled into a binary to avoid any
   *   lookup. Therefore, binary compatibility needs to be preserved
   *   on changes to types. (Use versioned type names to manage
   *   breaking changes.)
   * Note: this functionality is not currently available in the official
   * protobuf release, and it is not used for type URLs beginning with
   * type.googleapis.com.
   * Schemes other than `http`, `https` (or the empty scheme) might be
   * used with implementation specific semantics.
   */
  "@type"?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

export interface SlashrefundDeposit {
  depositor_address?: string;
  validator_address?: string;
  shares?: string;
}

/**
 * TODO: to account for more than one token, Tokens and Shares must be a struct.
 */
export interface SlashrefundDepositPool {
  operator_address?: string;

  /**
   * Coin defines a token with a denomination and an amount.
   *
   * NOTE: The amount field is an Int which implements the custom method
   * signatures required by gogoproto.
   */
  tokens?: V1Beta1Coin;
  shares?: string;
}

export type SlashrefundMsgClaimResponse = object;

export type SlashrefundMsgDepositResponse = object;

export interface SlashrefundMsgWithdrawResponse {
  /** @format date-time */
  completion_time?: string;
}

/**
 * Params defines the parameters for the module.
 */
export interface SlashrefundParams {
  allowedTokens?: string[];
}

export interface SlashrefundQueryAllDepositPoolResponse {
  depositPool?: SlashrefundDepositPool[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

export interface SlashrefundQueryAllDepositResponse {
  deposit?: SlashrefundDeposit[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

export interface SlashrefundQueryAllRefundPoolResponse {
  refundPool?: SlashrefundRefundPool[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

export interface SlashrefundQueryAllRefundResponse {
  refund?: SlashrefundRefund[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

export interface SlashrefundQueryAllUnbondingDepositResponse {
  unbondingDeposit?: SlashrefundUnbondingDeposit[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

export interface SlashrefundQueryGetDepositPoolResponse {
  /** TODO: to account for more than one token, Tokens and Shares must be a struct. */
  depositPool?: SlashrefundDepositPool;
}

export interface SlashrefundQueryGetDepositResponse {
  deposit?: SlashrefundDeposit;
}

export interface SlashrefundQueryGetRefundPoolResponse {
  /** TODO: to account for more than one token, Tokens and Shares must be a struct. */
  refundPool?: SlashrefundRefundPool;
}

export interface SlashrefundQueryGetRefundResponse {
  refund?: SlashrefundRefund;
}

export interface SlashrefundQueryGetUnbondingDepositResponse {
  unbondingDeposit?: SlashrefundUnbondingDeposit;
}

/**
 * QueryParamsResponse is response type for the Query/Params RPC method.
 */
export interface SlashrefundQueryParamsResponse {
  /** params holds all the parameters of this module. */
  params?: SlashrefundParams;
}

export interface SlashrefundRefund {
  delegator_address?: string;
  validator_address?: string;
  shares?: string;
}

/**
 * TODO: to account for more than one token, Tokens and Shares must be a struct.
 */
export interface SlashrefundRefundPool {
  operator_address?: string;

  /**
   * Coin defines a token with a denomination and an amount.
   *
   * NOTE: The amount field is an Int which implements the custom method
   * signatures required by gogoproto.
   */
  tokens?: V1Beta1Coin;
  shares?: string;
}

export interface SlashrefundUnbondingDeposit {
  depositorAddress?: string;
  validatorAddress?: string;
  entries?: SlashrefundUnbondingDepositEntry[];
}

export interface SlashrefundUnbondingDepositEntry {
  /** @format int64 */
  creation_height?: string;

  /** @format date-time */
  completion_time?: string;
  initial_balance?: string;
  balance?: string;
}

/**
* Coin defines a token with a denomination and an amount.

NOTE: The amount field is an Int which implements the custom method
signatures required by gogoproto.
*/
export interface V1Beta1Coin {
  denom?: string;
  amount?: string;
}

/**
* message SomeRequest {
         Foo some_parameter = 1;
         PageRequest pagination = 2;
 }
*/
export interface V1Beta1PageRequest {
  /**
   * key is a value returned in PageResponse.next_key to begin
   * querying the next page most efficiently. Only one of offset or key
   * should be set.
   * @format byte
   */
  key?: string;

  /**
   * offset is a numeric offset that can be used when key is unavailable.
   * It is less efficient than using key. Only one of offset or key should
   * be set.
   * @format uint64
   */
  offset?: string;

  /**
   * limit is the total number of results to be returned in the result page.
   * If left empty it will default to a value to be set by each app.
   * @format uint64
   */
  limit?: string;

  /**
   * count_total is set to true  to indicate that the result set should include
   * a count of the total number of items available for pagination in UIs.
   * count_total is only respected when offset is used. It is ignored when key
   * is set.
   */
  count_total?: boolean;

  /**
   * reverse is set to true if results are to be returned in the descending order.
   *
   * Since: cosmos-sdk 0.43
   */
  reverse?: boolean;
}

/**
* PageResponse is to be embedded in gRPC response messages where the
corresponding request message has used PageRequest.

 message SomeResponse {
         repeated Bar results = 1;
         PageResponse page = 2;
 }
*/
export interface V1Beta1PageResponse {
  /**
   * next_key is the key to be passed to PageRequest.key to
   * query the next page most efficiently. It will be empty if
   * there are no more results.
   * @format byte
   */
  next_key?: string;

  /**
   * total is total number of results available if PageRequest.count_total
   * was set, its value is undefined otherwise
   * @format uint64
   */
  total?: string;
}

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, ResponseType } from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({ securityWorker, secure, format, ...axiosConfig }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "" });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  private mergeRequestParams(params1: AxiosRequestConfig, params2?: AxiosRequestConfig): AxiosRequestConfig {
    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.instance.defaults.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  private createFormData(input: Record<string, unknown>): FormData {
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      formData.append(
        key,
        property instanceof Blob
          ? property
          : typeof property === "object" && property !== null
          ? JSON.stringify(property)
          : `${property}`,
      );
      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = (format && this.format) || void 0;

    if (type === ContentType.FormData && body && body !== null && typeof body === "object") {
      requestParams.headers.common = { Accept: "*/*" };
      requestParams.headers.post = {};
      requestParams.headers.put = {};

      body = this.createFormData(body as Record<string, unknown>);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        ...(requestParams.headers || {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title slashrefund/deposit.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryDepositAll
   * @summary Queries a list of Deposit items.
   * @request GET:/made-in-block/slash-refund/slashrefund/deposit
   */
  queryDepositAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<SlashrefundQueryAllDepositResponse, RpcStatus>({
      path: `/made-in-block/slash-refund/slashrefund/deposit`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryDeposit
   * @summary Queries a Deposit by index.
   * @request GET:/made-in-block/slash-refund/slashrefund/deposit/{depositorAddress}/{validatorAddress}
   */
  queryDeposit = (depositorAddress: string, validatorAddress: string, params: RequestParams = {}) =>
    this.request<SlashrefundQueryGetDepositResponse, RpcStatus>({
      path: `/made-in-block/slash-refund/slashrefund/deposit/${depositorAddress}/${validatorAddress}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryDepositPoolAll
   * @summary Queries a list of DepositPool items.
   * @request GET:/made-in-block/slash-refund/slashrefund/deposit_pool
   */
  queryDepositPoolAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<SlashrefundQueryAllDepositPoolResponse, RpcStatus>({
      path: `/made-in-block/slash-refund/slashrefund/deposit_pool`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryDepositPool
   * @summary Queries a DepositPool by index.
   * @request GET:/made-in-block/slash-refund/slashrefund/deposit_pool/{operatorAddress}
   */
  queryDepositPool = (operatorAddress: string, params: RequestParams = {}) =>
    this.request<SlashrefundQueryGetDepositPoolResponse, RpcStatus>({
      path: `/made-in-block/slash-refund/slashrefund/deposit_pool/${operatorAddress}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryParams
   * @summary Parameters queries the parameters of the module.
   * @request GET:/made-in-block/slash-refund/slashrefund/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<SlashrefundQueryParamsResponse, RpcStatus>({
      path: `/made-in-block/slash-refund/slashrefund/params`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryRefundAll
   * @summary Queries a list of Refund items.
   * @request GET:/made-in-block/slash-refund/slashrefund/refund
   */
  queryRefundAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<SlashrefundQueryAllRefundResponse, RpcStatus>({
      path: `/made-in-block/slash-refund/slashrefund/refund`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryRefund
   * @summary Queries a Refund by index.
   * @request GET:/made-in-block/slash-refund/slashrefund/refund/{delegator}/{validator}
   */
  queryRefund = (delegator: string, validator: string, params: RequestParams = {}) =>
    this.request<SlashrefundQueryGetRefundResponse, RpcStatus>({
      path: `/made-in-block/slash-refund/slashrefund/refund/${delegator}/${validator}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryRefundPoolAll
   * @summary Queries a list of RefundPool items.
   * @request GET:/made-in-block/slash-refund/slashrefund/refund_pool
   */
  queryRefundPoolAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<SlashrefundQueryAllRefundPoolResponse, RpcStatus>({
      path: `/made-in-block/slash-refund/slashrefund/refund_pool`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryRefundPool
   * @summary Queries a RefundPool by index.
   * @request GET:/made-in-block/slash-refund/slashrefund/refund_pool/{operatorAddress}
   */
  queryRefundPool = (operatorAddress: string, params: RequestParams = {}) =>
    this.request<SlashrefundQueryGetRefundPoolResponse, RpcStatus>({
      path: `/made-in-block/slash-refund/slashrefund/refund_pool/${operatorAddress}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryUnbondingDepositAll
   * @summary Queries a list of UnbondingDeposit items.
   * @request GET:/made-in-block/slash-refund/slashrefund/unbonding_deposit
   */
  queryUnbondingDepositAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<SlashrefundQueryAllUnbondingDepositResponse, RpcStatus>({
      path: `/made-in-block/slash-refund/slashrefund/unbonding_deposit`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryUnbondingDeposit
   * @summary Queries a UnbondingDeposit by index.
   * @request GET:/made-in-block/slash-refund/slashrefund/unbonding_deposit/{depositorAddress}/{validatorAddress}
   */
  queryUnbondingDeposit = (depositorAddress: string, validatorAddress: string, params: RequestParams = {}) =>
    this.request<SlashrefundQueryGetUnbondingDepositResponse, RpcStatus>({
      path: `/made-in-block/slash-refund/slashrefund/unbonding_deposit/${depositorAddress}/${validatorAddress}`,
      method: "GET",
      format: "json",
      ...params,
    });
}
