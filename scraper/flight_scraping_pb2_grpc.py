# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import flight_scraping_pb2 as flight__scraping__pb2


class FlightScraperStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetSouthwestHeaders = channel.unary_unary(
                '/flights.FlightScraper/GetSouthwestHeaders',
                request_serializer=flight__scraping__pb2.Empty.SerializeToString,
                response_deserializer=flight__scraping__pb2.SouthwestHeadersResponse.FromString,
                )


class FlightScraperServicer(object):
    """Missing associated documentation comment in .proto file."""

    def GetSouthwestHeaders(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_FlightScraperServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetSouthwestHeaders': grpc.unary_unary_rpc_method_handler(
                    servicer.GetSouthwestHeaders,
                    request_deserializer=flight__scraping__pb2.Empty.FromString,
                    response_serializer=flight__scraping__pb2.SouthwestHeadersResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'flights.FlightScraper', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class FlightScraper(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def GetSouthwestHeaders(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/flights.FlightScraper/GetSouthwestHeaders',
            flight__scraping__pb2.Empty.SerializeToString,
            flight__scraping__pb2.SouthwestHeadersResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
