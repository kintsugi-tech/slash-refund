/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

export interface DVPairs {
  dVPair: string;
}

const baseDVPairs: object = { dVPair: "" };

export const DVPairs = {
  encode(message: DVPairs, writer: Writer = Writer.create()): Writer {
    if (message.dVPair !== "") {
      writer.uint32(10).string(message.dVPair);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): DVPairs {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseDVPairs } as DVPairs;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.dVPair = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DVPairs {
    const message = { ...baseDVPairs } as DVPairs;
    if (object.dVPair !== undefined && object.dVPair !== null) {
      message.dVPair = String(object.dVPair);
    } else {
      message.dVPair = "";
    }
    return message;
  },

  toJSON(message: DVPairs): unknown {
    const obj: any = {};
    message.dVPair !== undefined && (obj.dVPair = message.dVPair);
    return obj;
  },

  fromPartial(object: DeepPartial<DVPairs>): DVPairs {
    const message = { ...baseDVPairs } as DVPairs;
    if (object.dVPair !== undefined && object.dVPair !== null) {
      message.dVPair = object.dVPair;
    } else {
      message.dVPair = "";
    }
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;
