defmodule JankoMuzykant.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    children = [
      JankoMuzykantWeb.Telemetry,
      JankoMuzykant.Repo,
      {DNSCluster, query: Application.get_env(:janko_muzykant, :dns_cluster_query) || :ignore},
      {Phoenix.PubSub, name: JankoMuzykant.PubSub},
      # Start the Finch HTTP client for sending emails
      {Finch, name: JankoMuzykant.Finch},
      # Start a worker by calling: JankoMuzykant.Worker.start_link(arg)
      # {JankoMuzykant.Worker, arg},
      # Start to serve requests, typically the last entry
      JankoMuzykantWeb.Endpoint
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: JankoMuzykant.Supervisor]
    Supervisor.start_link(children, opts)
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  @impl true
  def config_change(changed, _new, removed) do
    JankoMuzykantWeb.Endpoint.config_change(changed, removed)
    :ok
  end
end
