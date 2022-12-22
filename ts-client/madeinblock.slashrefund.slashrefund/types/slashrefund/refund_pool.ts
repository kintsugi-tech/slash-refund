/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Coin } from "../cosmos/base/v1beta1/coin";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

/** TODO: to account for more than one token, Tokens and Shares must be a struct. */
export interface RefundPool {
  operatorAddress: string;
  tokens: Coin | undefined;
  shares: string;
}

function createBaseRefundPool(): RefundPool {
  return { operatorAddress: "", tokens: undefined, shares: "" };
}

export const RefundPool = {
  encode(message: RefundPool, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.operatorAddress !== "") {
      writer.uint32(10).string(message.operatorAddress);
    }
    if (message.tokens !== undefined) {
      Coin.encode(message.tokens, writer.uint32(18).fork()).ldelim();
    }
    if (message.shares !== "") {
      writer.uint32(26).string(message.shares);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RefundPool {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRefundPool();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.operatorAddress = reader.string();
          break;
        case 2:
          message.tokens = Coin.decode(reader, reader.uint32());
          break;
        case 3:
          message.shares = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RefundPool {
    return {
      operatorAddress: isSet(object.operatorAddress) ? String(object.operatorAddress) : "",
      tokens: isSet(object.tokens) ? Coin.fromJSON(object.tokens) : undefined,
      shares: isSet(object.shares) ? String(object.shares) : "",
    };
  },

  toJSON(message: RefundPool): unknown {
    const obj: any = {};
    message.operatorAddress !== undefined && (obj.operatorAddress = message.operatorAddress);
    message.tokens !== undefined && (obj.tokens = message.tokens ? Coin.toJSON(message.tokens) : undefined);
    message.shares !== undefined && (obj.shares = message.shares);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<RefundPool>, I>>(object: I): RefundPool {
    const message = createBaseRefundPool();
    message.operatorAddress = object.operatorAddress ?? "";
    message.tokens = (object.tokens !== undefined && object.tokens !== null)
      ? Coin.fromPartial(object.tokens)
      : undefined;
    message.shares = object.shares ?? "";
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
