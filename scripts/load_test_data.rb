#!/usr/bin/env ruby
# frozen_string_literal: true

require 'net/http'
require 'json'
require 'uri'
require 'optparse'

class MeilisearchLoader
  attr_reader :host, :master_key, :index_name, :verbose

  def initialize(options = {})
    @host = options[:host] || 'http://localhost:7700'
    @master_key = options[:master_key] || 'MASTER_KEY'
    @index_name = options[:index_name] || 'movies'
    @verbose = options[:verbose] || false
  end

  def run
    create_index
    configure_settings
    add_documents
    configure_synonyms
    run_test_search
  end

  private

  def log(message)
    puts message if verbose
  end

  def request(method, path, data = nil)
    uri = URI.parse("#{host}#{path}")
    http = Net::HTTP.new(uri.host, uri.port)

    case method.to_s.downcase
    when 'get'
      request = Net::HTTP::Get.new(uri.request_uri)
    when 'post'
      request = Net::HTTP::Post.new(uri.request_uri)
    when 'put'
      request = Net::HTTP::Put.new(uri.request_uri)
    when 'delete'
      request = Net::HTTP::Delete.new(uri.request_uri)
    else
      raise "Unsupported HTTP method: #{method}"
    end

    request['Content-Type'] = 'application/json'
    request['Authorization'] = "Bearer #{master_key}"
    request.body = data.to_json if data

    log "#{method.upcase} #{uri}"
    log "Request body: #{data.to_json}" if data && verbose

    response = http.request(request)

    log "Response status: #{response.code}"
    log "Response body: #{response.body}" if verbose

    return response
  end

  def create_index
    log "Creating index: #{index_name}"
    response = request(:post, '/indexes', {
      uid: index_name,
      primaryKey: 'id'
    })

    case response.code.to_i
    when 201
      log "Index created successfully"
    when 400
      log "Index creation failed: #{response.body}"
    when 409
      log "Index already exists"
    else
      log "Unexpected response: #{response.code} - #{response.body}"
    end
  end

  def configure_settings
    log "Configuring searchable attributes"
    request(:put, "/indexes/#{index_name}/settings/searchable-attributes",
      ['title', 'overview', 'genres'])

    log "Configuring filterable attributes"
    request(:put, "/indexes/#{index_name}/settings/filterable-attributes",
      ['genres', 'release_date', 'rating'])
  end

  def add_documents
    log "Adding documents to index"
    documents = [
      {
        id: 1,
        title: "Carol",
        overview: "In 1950s New York, a department-store clerk who dreams of a better life falls for an older, married woman.",
        genres: ["Romance", "Drama"],
        release_date: 1448582400,
        rating: 7.2
      },
      {
        id: 2,
        title: "Wonder Woman",
        overview: "An Amazon princess comes to the world of Man in the grips of the First World War to confront the forces of evil and bring an end to human conflict.",
        genres: ["Action", "Adventure", "Fantasy"],
        release_date: 1496361600,
        rating: 7.5
      },
      {
        id: 3,
        title: "Life of Pi",
        overview: "After a shipwreck, a young man adrift in the ocean aboard a lifeboat shares his small craft with an adult Bengal tiger.",
        genres: ["Adventure", "Drama", "Fantasy"],
        release_date: 1353542400,
        rating: 7.9
      },
      {
        id: 4,
        title: "Mad Max: Fury Road",
        overview: "In a post-apocalyptic wasteland, a woman rebels against a tyrannical ruler in search for her homeland with the aid of a group of female prisoners, a psychotic worshiper, and a drifter named Max.",
        genres: ["Action", "Adventure", "Science Fiction"],
        release_date: 1431648000,
        rating: 8.1
      },
      {
        id: 5,
        title: "Moana",
        overview: "In Ancient Polynesia, when a terrible curse incurred by the Demigod Maui reaches Moana's island, she answers the Ocean's call to seek out the Demigod to set things right.",
        genres: ["Animation", "Family", "Adventure"],
        release_date: 1479859200,
        rating: 7.6
      },
      {
        id: 6,
        title: "Philadelphia",
        overview: "Two competing lawyers join forces to sue a prestigious law firm for AIDS discrimination. As their unlikely friendship develops, their courage overcomes the prejudice and corruption of their powerful adversaries.",
        genres: ["Drama"],
        release_date: 758505600,
        rating: 7.7
      },
      {
        id: 7,
        title: "Inception",
        overview: "A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O.",
        genres: ["Action", "Science Fiction", "Adventure"],
        release_date: 1279238400,
        rating: 8.3
      },
      {
        id: 8,
        title: "The Shawshank Redemption",
        overview: "Framed in the 1940s for the double murder of his wife and her lover, upstanding banker Andy Dufresne begins a new life at the Shawshank prison, where he puts his accounting skills to work for an amoral warden.",
        genres: ["Drama", "Crime"],
        release_date: 780278400,
        rating: 8.7
      },
      {
        id: 9,
        title: "Parasite",
        overview: "All unemployed, Ki-taek's family takes peculiar interest in the wealthy and glamorous Parks for their livelihood until they get entangled in an unexpected incident.",
        genres: ["Comedy", "Thriller", "Drama"],
        release_date: 1557964800,
        rating: 8.5
      },
      {
        id: 10,
        title: "Your Name",
        overview: "High schoolers Mitsuha and Taki are complete strangers living separate lives. But one night, they suddenly switch places. Mitsuha wakes up in Taki's body, and he in hers. This bizarre occurrence continues to happen randomly, and the two must adjust their lives around each other.",
        genres: ["Romance", "Animation", "Drama"],
        release_date: 1470960000,
        rating: 8.4
      }
    ]

    response = request(:post, "/indexes/#{index_name}/documents", documents)

    if response.code.to_i == 202
      task = JSON.parse(response.body)
      log "Documents added successfully. Task ID: #{task['taskUid']}"
    else
      log "Failed to add documents: #{response.body}"
    end
  end

  def configure_synonyms
    log "Configuring synonyms"
    synonyms = {
      "great" => ["fantastic", "excellent"],
      "fantastic" => ["great", "excellent"],
      "sci-fi" => ["science fiction"],
      "science fiction" => ["sci-fi"]
    }

    response = request(:put, "/indexes/#{index_name}/settings/synonyms", synonyms)

    if response.code.to_i == 202
      log "Synonyms configured successfully"
    else
      log "Failed to configure synonyms: #{response.body}"
    end
  end

  def run_test_search
    log "Running test search for 'adventure'"
    response = request(:post, "/indexes/#{index_name}/search", { q: "adventure" })

    if response.code.to_i == 200
      result = JSON.parse(response.body)
      log "Search returned #{result['hits'].size} results"
      result['hits'].each do |hit|
        puts "- #{hit['title']} (#{hit['rating']})"
      end
    else
      log "Search failed: #{response.body}"
    end

    log "\nRunning filtered search for 'adventure' with rating > 7.8"
    response = request(:post, "/indexes/#{index_name}/search", {
      q: "adventure",
      filter: "rating > 7.8"
    })

    if response.code.to_i == 200
      result = JSON.parse(response.body)
      log "Filtered search returned #{result['hits'].size} results"
      result['hits'].each do |hit|
        puts "- #{hit['title']} (#{hit['rating']})"
      end
    else
      log "Filtered search failed: #{response.body}"
    end
  end
end

# Parse command line options
options = {}
OptionParser.new do |opts|
  opts.banner = "Usage: load_test_data.rb [options]"

  opts.on("-h", "--host HOST", "Meilisearch host URL (default: http://localhost:7700)") do |h|
    options[:host] = h
  end

  opts.on("-k", "--key KEY", "Master key (default: MASTER_KEY)") do |k|
    options[:master_key] = k
  end

  opts.on("-i", "--index NAME", "Index name (default: movies)") do |i|
    options[:index_name] = i
  end

  opts.on("-v", "--verbose", "Enable verbose output") do |v|
    options[:verbose] = v
  end

  opts.on("--help", "Show this help message") do
    puts opts
    exit
  end
end.parse!

# Run the loader
loader = MeilisearchLoader.new(options)
loader.run
