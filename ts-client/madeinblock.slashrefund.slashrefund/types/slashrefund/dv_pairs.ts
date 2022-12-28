/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { DVPair } from "./dv_pair";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

export interface DVPairs {
  pairs: DVPair[];
}

function createBaseDVPairs(): DVPairs {
  return { pairs: [] };
}

export const DVPairs = {
  encode(message: DVPairs, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.pairs) {
      DVPair.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DVPairs {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDVPairs();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pairs.push(DVPair.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DVPairs {
    return { pairs: Array.isArray(object?.pairs) ? object.pairs.map((e: any) => DVPair.fromJSON(e)) : [] };
  },

  toJSON(message: DVPairs): unknown {
    const obj: any = {};
    if (message.pairs) {
      obj.pairs = message.pairs.map((e) => e ? DVPair.toJSON(e) : undefined);
    } else {
      obj.pairs = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DVPairs>, I>>(object: I): DVPairs {
    const message = createBaseDVPairs();
    message.pairs = object.pairs?.map((e) => DVPair.fromPartial(e)) || [];
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
