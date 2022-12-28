/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

export interface Deposit {
  depositorAddress: string;
  validatorAddress: string;
  shares: string;
}

function createBaseDeposit(): Deposit {
  return { depositorAddress: "", validatorAddress: "", shares: "" };
}

export const Deposit = {
  encode(message: Deposit, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.depositorAddress !== "") {
      writer.uint32(10).string(message.depositorAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    if (message.shares !== "") {
      writer.uint32(26).string(message.shares);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Deposit {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeposit();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.depositorAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
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

  fromJSON(object: any): Deposit {
    return {
      depositorAddress: isSet(object.depositorAddress) ? String(object.depositorAddress) : "",
      validatorAddress: isSet(object.validatorAddress) ? String(object.validatorAddress) : "",
      shares: isSet(object.shares) ? String(object.shares) : "",
    };
  },

  toJSON(message: Deposit): unknown {
    const obj: any = {};
    message.depositorAddress !== undefined && (obj.depositorAddress = message.depositorAddress);
    message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
    message.shares !== undefined && (obj.shares = message.shares);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Deposit>, I>>(object: I): Deposit {
    const message = createBaseDeposit();
    message.depositorAddress = object.depositorAddress ?? "";
    message.validatorAddress = object.validatorAddress ?? "";
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
